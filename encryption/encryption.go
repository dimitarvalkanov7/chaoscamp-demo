package encryption

import (
	"crypto/aes"
	"encoding/hex"
	"log"
)

var (
	key = []byte("ZxxULO156stVFP3UzfcLM1Dn32EESlePg44")
)

func Encrypt(plaintext string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("NewCipher(%d bytes) = %s", len(key), err)
		panic(err)
	}
	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func Decrypt(ct string) {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("NewCipher(%d bytes) = %s", len(key), err)
		panic(err)
	}
	plain := make([]byte, len(ciphertext))
	c.Decrypt(plain, ciphertext)
	s := string(plain[:])
	log.Printf("AES Decrypyed Text:  %s\n", s)
}
