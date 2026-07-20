package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/routes"
)

func main() {

	config.LoadEnv()

	config.ConnectDatabase() // konek database
	config.ConnectRedis()    // konek redis

	router := gin.Default()

	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "API Running",
	// 	})
	// })

	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "API Running",
	// 	})
	// })

	routes.SetupRoutes(router)

	router.Run(":8080")
}
