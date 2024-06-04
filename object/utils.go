package object

import (
	"crypto/sha1"
	"encoding/hex"
)

// Pattern from https://play.golang.org/p/YUaWWEeB4U
func generateSHA1Hash(data string) []byte {
	hash := sha1.New()
	hash.Write([]byte(data))
	byteSlice := hash.Sum(nil)

	return byteSlice
}

// Note: to convert a "sha1 value to string" it needs to be encoded first, since a hash is binary.
// The traditional encoding for SHA hashes is hex (import "encoding/hex"):
func GetSha1AsString(hash []byte) string {
	return hex.EncodeToString(hash)
}
