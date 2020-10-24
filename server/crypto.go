package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func decrypt(aesKey, iv, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func encrypt(plaintext []byte) (aesKey, iv, ciphertext []byte, err error) {
	aesKey = make([]byte, 32)
	_, err = rand.Read(aesKey)
	if err != nil {
		return nil, nil, nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, nil, err
	}
	iv = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	return aesKey, iv, aesgcm.Seal(nil, iv, plaintext, nil), nil
}

