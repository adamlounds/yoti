package crypto

import (
	"encoding/hex"
	"testing"
)

func TestEncrypt(t *testing.T) {
	aesKey, body, err := Encrypt([]byte(""))
	if err != nil {
		t.Errorf("Gave error on simple encrypt")
	}
	if len(aesKey) != 32 {
		t.Errorf("aesKey length incorrect, got %d, want %d", len(aesKey), 32)
	}
	if len(body) != 28 {
		t.Errorf("body length incorrect, got %d, want %d", len(body), 28)
	}

	aesKey, body, err = Encrypt([]byte("plaintext"))
	if err != nil {
		t.Errorf("Gave error on simple encrypt")
	}
	if len(aesKey) != 32 {
		t.Errorf("aesKey length incorrect, got %d, want %d", len(aesKey), 32)
	}
	if len(body) != 37 {
		t.Errorf("body length incorrect, got %d, want %d", len(body), 37)
	}
}

func TestDecrypt(t *testing.T) {
	aesKey, _ := hex.DecodeString("eaf2f9032287afb5c480e72d48d3a82000a787c80f11760e7671f72b44d303d4")
	body, _ := hex.DecodeString("6d495228ca04993c07b1bc89dbcdde0535647e013d05a1bf71694905e098632853b6fc1678")

	plaintext, err := Decrypt(aesKey, body)
	if err != nil {
		t.Errorf("error on decrypt")
	}
	if string(plaintext) != "plaintext" {
		t.Errorf("invalid decrypt, got %v", plaintext)
	}

	body[0] = 'a'
	plaintext, err = Decrypt(aesKey, body)
	if err == nil {
		t.Errorf("missing error when decrypting bad data")
	}
}
