package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env not found, using system environment")
	}
	return nil
}
