package main

import(
	"fmt"
	"github.com/adamlounds/yoti/client"
	"github.com/rs/zerolog"
	"os"
	"time"
)

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



}