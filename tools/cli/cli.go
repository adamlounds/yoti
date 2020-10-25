package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/adamlounds/yoti/client"
)

const (
	DefaultEndpoint = "http://localhost:8080"
)

func main() {
	endpoint := flag.String("endpoint", DefaultEndpoint, "server endpoint")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify action: store with id and key or retrieve with filename\n")
		os.Exit(1)
	}
	action := flag.Arg(0)

	if action == "store" && flag.NArg() != 3 {
		fmt.Fprintf(os.Stderr, "ERROR: you need to specify <id> and <filename>\n")
		os.Exit(1)
	}
	if action == "retrieve" && flag.NArg() != 3 {
		fmt.Fprintf(os.Stderr, "ERROR: you need to specify <id> and <key> to retrieve\n")
		os.Exit(1)
	}

	client, err := client.NewClient(&client.Config{
		Endpoint: *endpoint,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot create client (%s)", err);
		os.Exit(1)
	}

	if action == "store" {
		id := []byte(flag.Arg(1))
		payload, err := ioutil.ReadFile(flag.Arg(2))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot read file (%s)", err)
			os.Exit(1)
		}
		aesKey, err := client.Store(id, payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot store (%s)", err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, string(aesKey), "\n")
		os.Exit(0)
	}
	if action == "retrieve" {
		id := []byte(flag.Arg(1))

		aesKey := make([]byte, 32)
		n, err := hex.Decode(aesKey, []byte(flag.Arg(2)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: bad aesKey (%s)", err)
			os.Exit(1)
		}
		if n != 32 {
			fmt.Fprintf(os.Stderr, "ERROR: aesKey must be 64 hexits", err)
			os.Exit(1)
		}

		payload, err := client.Retrieve(id, aesKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot retrieve (%s)", err)
			os.Exit(1)
		}

		os.Stdout.Write(payload)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "ERROR: unknown command")
	os.Exit(1)
}
