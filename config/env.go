package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system environment")
	}
}

func Env(key string) string {
	return os.Getenv(key)
}