package pkg

import (
	"os"
)

func validateToken(token string) bool {
	if token == os.Getenv("API_KEY") {
		return true
	} else {
		return false
	}
}
