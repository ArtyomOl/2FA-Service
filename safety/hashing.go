package safety

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hashing(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	hashedData := hash.Sum(nil)
	hexHash := hex.EncodeToString(hashedData)
	return hexHash
}
