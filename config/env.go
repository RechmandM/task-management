package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// coba cari .env di folder saat ini
	if err := godotenv.Load(".env"); err == nil {
		return
	}

	// kalau test dijalankan dari folder tests
	if err := godotenv.Load("../.env"); err == nil {
		return
	}

	log.Println(".env not found, using system environment")
}

func Env(key string) string {
	return os.Getenv(key)
}
