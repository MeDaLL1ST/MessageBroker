package pkg

import (
	"github.com/joho/godotenv"
	"log"
)

func Loadenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
