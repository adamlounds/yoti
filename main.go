package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adamlounds/yoti/crypto"
	"github.com/adamlounds/yoti/dao"
	"github.com/adamlounds/yoti/store"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

type StoreRequest struct {
	Id      []byte
	Payload []byte
}
type RetrieveRequest struct {
	Id     []byte
	AesKey []byte
}

const secretSalt = "1911797e2e9d418b8399fafd79de79f14c6370ae58c2a314195a35bcfdd359ae"

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	storeFac, err := store.NewStoreFactory()
	if err != nil {
		logger.Error().Err(err).Msg("Cannot create data stores")
		os.Exit(1)
	}

	http.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		reqId := ulid.MustNew(ulid.Now(), rand.Reader)
		reqLogger := logger.With().Str("reqid", reqId.String()).Logger()
		daoFac := dao.NewFactory(reqLogger, storeFac)

		// TODO move to middleware, add more details, probably use a standard
		// per-request logging middleware
		defer func(begin time.Time) {
			reqLogger.Info().
				Dur("duration", time.Since(begin)).
				Msg("store")
		}(time.Now())
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reqLogger.Warn().Err(err).Msg("cannot read body")
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var req StoreRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			reqLogger.Warn().Err(err).Msg("cannot unmarshal json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		aesKey, ciphertext, err := crypto.Encrypt(req.Payload)
		if err != nil {
			reqLogger.Warn().Err(err).Msg("cannot encrypt")
			http.Error(w, "cannot read body", http.StatusInternalServerError)
			return
		}

		idSalt, _ := hex.DecodeString(secretSalt)
		saltedId := append(req.Id, idSalt...)
		storedId := sha256.Sum256(saltedId)

		err = daoFac.Document.Store(storedId[:], ciphertext)
		if err != nil {
			reqLogger.Warn().Err(err).Msg("cannot store")
			http.Error(w, "cannot store", http.StatusInternalServerError)
			return
		}

		hexKey := make([]byte, hex.EncodedLen(len(aesKey)))
		_ = hex.Encode(hexKey, aesKey)
		_, _ = w.Write(hexKey) // nocheck: if write fails then any error is also likely to fail
	})

	http.HandleFunc("/retrieve", func(w http.ResponseWriter, r *http.Request) {
		reqId := ulid.MustNew(ulid.Now(), rand.Reader)
		reqLogger := logger.With().Str("reqid", reqId.String()).Logger()
		daoFac := dao.NewFactory(reqLogger, storeFac)

		defer func(begin time.Time) {
			reqLogger.Info().
				Dur("duration", time.Since(begin)).
				Msg("retrieve")
		}(time.Now())
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reqLogger.Warn().Msg("cannot read body")
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var req RetrieveRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			reqLogger.Warn().Err(err).Msg("cannot unmarshal json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if len(req.AesKey) != 32 {
			reqLogger.Warn().Bytes("aesKey", req.AesKey).Msg("short/long aesKey")
			http.Error(w, "short/long key", http.StatusBadRequest)
			return
		}
		idSalt, _ := hex.DecodeString(secretSalt)
		saltedId := append(req.Id, idSalt...)
		storedId := sha256.Sum256(saltedId)

		ciphertext, err := daoFac.Document.Retrieve(storedId[:])
		if err != nil {
			reqLogger.Warn().Bytes("id", req.Id).Err(err).Msg("cannot fetch")
			http.Error(w, "cannot fetch", http.StatusInternalServerError)
			return
		}
		payload, err := crypto.Decrypt(req.AesKey, ciphertext)
		if err != nil {
			reqLogger.Warn().Bytes("id", req.Id).Err(err).Msg("cannot decrypt")
			http.Error(w, "cannot decrypt", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(payload)
	})

	logger.Info().Int("port", 8080).Msg("starting")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
