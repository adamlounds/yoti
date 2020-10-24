package server

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Entry struct {
	ciphertext []byte
	salt []byte
}

type ClientInstance struct {
	DataStore map[string]Entry
}
const secretSalt = "1911797e2e9d418b8399fafd79de79f14c6370ae58c2a314195a35bcfdd359ae"

func (c ClientInstance) Store(id, payload []byte) (aesKey []byte, err error) {
	aesKey, salt, ciphertext, err := encrypt(payload)
	if err != nil {
		//fmt.Printf("cannot encrypt! %v", err)
		return nil, err
	}

	idSalt, _ := hex.DecodeString(secretSalt)
	saltedId := append(id, idSalt...)
	storedId := sha1.Sum(saltedId)
	c.DataStore[string(storedId[:])] = Entry{ciphertext, salt} // TODO mutex/safe update
	return aesKey, nil
}

func (c ClientInstance) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	idSalt, _ := hex.DecodeString(secretSalt)
	saltedId := append(id, idSalt...);
	storedId := sha1.Sum(saltedId)
	entry, ok := c.DataStore[string(storedId[:])]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}
	//fmt.Printf("ciphertext %s id %s ok %v\n", string(entry.ciphertext), storedId, ok)

	plaintext, err := decrypt(aesKey, entry.salt, entry.ciphertext)
	if err != nil {
		fmt.Printf("could not decrypt %v\n", plaintext)
		return nil, err
	}
	//fmt.Printf("decrypted to %s\n", string(plaintext))
	return plaintext, nil
}


