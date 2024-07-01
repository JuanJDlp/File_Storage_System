package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

//HashString creates a hash for the file name that will always be 64 bytes long
func HashString(fileName string) string {
	hash := sha256.New()
	hash.Write([]byte(fileName))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
