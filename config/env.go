package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Mencari file .env pada root project
	if err := godotenv.Load(".env"); err == nil {
		return
	}

	// Digunakan ketika unit test dijalankan dari folder tests
	if err := godotenv.Load("../.env"); err == nil {
		return
	}

	log.Println(".env not found, using system environment")
}

// Env digunakan untuk mengambil nilai environment variable.
func Env(key string) string {
	return os.Getenv(key)
}
