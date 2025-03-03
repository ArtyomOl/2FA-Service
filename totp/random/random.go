package random

import (
	"crypto/rand"
	"encoding/base32"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 10
)

func RandomToken() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

// GenerateRandomString генерирует строку случайных символов заданной длины
func RandomCode() string {
	result := make([]byte, length)
	for i := range result {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(61))
		n := nBig.Int64()
		result[i] = charset[n] // Выбираем случайный символ из charset
	}
	return string(result)
}
