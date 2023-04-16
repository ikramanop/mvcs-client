package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/ikramanop/mvcs-client/app/model"
)

const (
	GenerateRandomStringError = "ошибка при генерации строки"
	EncryptStringError        = "ошибка при шифровании строки"
	DecryptStringError        = "ошибка при дешифровании строки"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", model.WrapError(GenerateRandomStringError, err)
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)

	return encoded, nil
}

func EncryptString(str string, password string) (string, error) {
	bstr := []byte(str)
	bpass := []byte(password)

	blockCipher, err := aes.NewCipher(bpass)
	if err != nil {
		return "1", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "2", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "3", err
	}

	ciphertext := gcm.Seal(nonce, nonce, bstr, nil)

	return string(ciphertext), nil
}

func DecryptString(str string, password string) (string, error) {
	bstr := []byte(str)
	bpass := []byte(password)

	blockCipher, err := aes.NewCipher(bpass)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := bstr[:gcm.NonceSize()], bstr[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
