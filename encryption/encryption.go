package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	//"encoding/json"
	"errors"
	"github.com/dimitarvalkanov7/chaoscamp-demo/models"
	"net/http"
)

var (
	key = "123456789012345678901234"
	iv  = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
)

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return encodeBase64(ciphertext)
}

func Decrypt(text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext := decodeBase64(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}

func GetLoggedUser(r *http.Request) (*models.User, error) {
	c, err := r.Cookie("demoscenes")
	if err != nil {
		return nil, err
	}

	decrEmail := Decrypt(c.Value)

	var user *(models.User)
	user = user.GetUserByEmail(decrEmail)
	if user == nil {
		return nil, errors.New("No user with email: " + decrEmail)
	}

	return user, nil
}
