package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	Ctx   = context.Background()
)

func ConnectRedis() {

	db, _ := strconv.Atoi(Env("REDIS_DB"))

	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Env("REDIS_HOST"), Env("REDIS_PORT")),
		Password: Env("REDIS_PASS"),
		DB:       db,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("❌ Failed to connect Redis :", err)
	}

	log.Println("✅ Redis Connected")
}
