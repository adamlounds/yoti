package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Decrypt(aesKey, body []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv := body[:12]
	ciphertext := body[12:]
	plaintext, err = aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Encrypt(plaintext []byte) (aesKey, body []byte, err error) {
	aesKey = make([]byte, 32)
	_, err = rand.Read(aesKey)
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// need both iv+key to decrypt. Return key, store known-length iv with ciphertext.
	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)
	return aesKey, append(iv, ciphertext...), nil
}

