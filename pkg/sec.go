package pkg

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func generateChecksum(inhash, key, value string) bool {
	data := []byte(key + value + os.Getenv("API_KEY"))
	hash := sha256.Sum256(data)
	if inhash == fmt.Sprintf("%x", hash[:]) {
		return true
	} else {
		return false
	}
}

func validateToken(token string) bool {
	if token == os.Getenv("API_KEY") {
		return true
	} else {
		return false
	}
}
