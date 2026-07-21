package main

import (
	"time"

	"github.com/gin-contrib/cors" // 1. Import package cors
	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/routes"
)

func main() {

	// Load konfigurasi environment (.env)
	config.LoadEnv()

	// Inisialisasi koneksi database MySQL
	config.ConnectDatabase()

	// Inisialisasi koneksi Redis
	config.ConnectRedis()

	// Membuat instance Gin Engine
	router := gin.Default()

	// Middleware CORS
	// Mengizinkan frontend (React Native Web) mengakses API
	// selama proses development dan testing.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint sederhana untuk memastikan API berjalan
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Running",
		})
	})

	// Registrasi seluruh endpoint API
	routes.SetupRoutes(router)

	// Menjalankan server pada port 8080
	router.Run(":8080")
}
