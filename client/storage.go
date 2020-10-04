package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
)

type Entry struct {
	ciphertext []byte
	salt []byte
}

type ClientInstance struct {
	DataStore map[string]Entry
}

func (c ClientInstance) Store(id, payload []byte) (aesKey []byte, err error) {
	aesKey, salt, ciphertext, err := encrypt(payload)
	if err != nil {
		//fmt.Printf("cannot encrypt! %v", err)
		return nil, err
	}

	storedId := sha1.Sum(id)
	c.DataStore[string(storedId[:])] = Entry{ciphertext, salt} // TODO mutex/safe update
	return aesKey, nil
}

func (c ClientInstance) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	storedId := sha1.Sum(id)
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

func decrypt(aesKey, salt, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesgcm.Open(nil, salt, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func encrypt(plaintext []byte) (aesKey, salt, ciphertext []byte, err error) {
	aesKey = make([]byte, 32)
	_, err = rand.Read(aesKey)
	if err != nil {
		return nil, nil, nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, nil, err
	}
	salt = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, err
	}
	salt = []byte("SALTSALTSALT")

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	return aesKey, salt, aesgcm.Seal(nil, salt, plaintext, nil), nil

}

