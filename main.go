package main

import(
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/adamlounds/yoti/server"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type StoreRequest struct {
	Id string
	Payload string
}
type RetrieveRequest struct {
	Id string
	Key string
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	myClient := server.ClientInstance{DataStore: make(map[string]server.Entry)}

	http.HandleFunc("/store", func (w http.ResponseWriter, r * http.Request) {
		// TODO move to middleware, add more details, probably use a standard
		// per-request logging middleware
		defer func(begin time.Time) {
			logger.Info().
				Dur("duration", time.Since(begin)).
				Msg("store")
		}(time.Now())
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warn().Err(err).Msg("cannot read body")
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var req StoreRequest
		json.Unmarshal(body, &req)
		aesKey, err := myClient.Store([]byte(req.Id), []byte(req.Payload))
		if err != nil {
			logger.Warn().Err(err).Msg("cannot encrypt")
			http.Error(w, "cannot read body", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, hex.EncodeToString(aesKey))
	})

	http.HandleFunc("/retrieve", func (w http.ResponseWriter, r * http.Request) {
		defer func(begin time.Time) {
			logger.Info().
				Dur("duration", time.Since(begin)).
				Msg("retrieve")
		}(time.Now())
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warn().Msg("cannot read body")
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var req RetrieveRequest
		json.Unmarshal(body, &req)
		aesKey := make([]byte, 32)
		n, err := hex.Decode(aesKey, []byte(req.Key))
		if err != nil {
			logger.Warn().Str("key", req.Key).Err(err).Msg("bad key")
			http.Error(w, "bad key", http.StatusBadRequest)
			return
		}
		if n != 32 {
			logger.Warn().Str("key", req.Key).Msg("short/long key")
			http.Error(w, "short/long key", http.StatusBadRequest)
			return
		}
		payload, err := myClient.Retrieve([]byte(req.Id), aesKey)
		if err != nil {
			logger.Warn().Str("id", req.Id).Err(err).Msg("cannot fetch/decrypt")
			http.Error(w, "cannot fetch/decrypt", http.StatusInternalServerError)
			return

		}
		fmt.Fprint(w, string(payload))
	})

	logger.Info().Int("port", 8080).Msg("starting")
	logger.Info().Err( http.ListenAndServe(":8080", nil))
}
