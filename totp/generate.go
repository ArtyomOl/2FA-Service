package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"math"
	"time"
)

const (
	timeStep = 30
	digits   = 6
)

func CreateTimeBasedCode(secret string) int {
	decodedSecret, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		panic(err)
	}

	timeCounter := time.Now().Unix() / timeStep

	timeByteArray := make([]byte, 8)
	for i := 0; i < 8; i++ {
		timeByteArray[7-i] = byte(timeCounter >> (8 * i))
	}

	h := hmac.New(sha1.New, decodedSecret)
	h.Write(timeByteArray)
	hmacResult := h.Sum(nil)

	offset := hmacResult[len(hmacResult)-1] & 0xf
	truncatedHash := (int(hmacResult[offset]&0x7f) << 24) |
		(int(hmacResult[offset+1]&0xff) << 16) |
		(int(hmacResult[offset+2]&0xff) << 8) |
		(int(hmacResult[offset+3] & 0xff))

	otp := truncatedHash % int(math.Pow(10, float64(digits)))

	return otp
}
