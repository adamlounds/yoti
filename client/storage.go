package client

import (
	"crypto/sha1"
	"fmt"
)

type ClientInstance struct {
	DataStore map[string][]byte
}

func (c ClientInstance) Store(id, payload []byte) (aesKey []byte, err error) {

	ciphertext := payload
	storedId := sha1.Sum(id)
	c.DataStore[string(storedId[:])] = ciphertext // TODO mutex/safe update
	return []byte("aesKey"), nil
}

func (c ClientInstance) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	storedId := sha1.Sum(id)
	ciphertext, ok := c.DataStore[string(storedId[:])]
	fmt.Printf("ciphertext %s id %s ok %v\n", string(ciphertext), storedId, ok)
	return ciphertext, nil
}

