package main

import (
	"time"

	"github.com/gin-contrib/cors" // 1. Import package cors
	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/routes"
)

func main() {

	config.LoadEnv() // load .env

	config.ConnectDatabase() // konek database
	config.ConnectRedis()    // konek redis

	router := gin.Default()

	// 2. Pasang Middleware CORS
	// Pengaturan AllowAllOrigins cocok & cepat untuk kebutuhan tes/development
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Running",
		})
	})

	routes.SetupRoutes(router)

	router.Run(":8080")
}

