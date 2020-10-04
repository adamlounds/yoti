package main

import(
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/adamlounds/yoti/client"
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

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	logger.Info().Msg("starting");
	var myClient client.ClientInstance
	myClient = client.ClientInstance{DataStore: make(map[string]client.Entry)}

	aesKey, _ := myClient.Store([]byte("abc"), []byte("plaintextxx"))

	payload, _ := myClient.Retrieve([]byte("abc"), aesKey)
	fmt.Printf("retrieved, got payload %s\n", string(payload))

	//payload, err := myClient.Retrieve([]byte("def"), aesKey)

	http.HandleFunc("/store", func (w http.ResponseWriter, r * http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warn().Msg("cannot read body")
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}

		var req StoreRequest
		json.Unmarshal(body, &req)
		aesKey, err := myClient.Store([]byte(req.Id), []byte(req.Payload))
		fmt.Fprint(w, hex.EncodeToString(aesKey))
	})

	logger.Info().Err( http.ListenAndServe(":8080", nil))
}
