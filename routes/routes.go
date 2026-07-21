package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/controllers"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api")
	{

		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.GetTasks)
			tasks.POST("", controllers.CreateTask)
			tasks.GET("/:id", controllers.GetTask)
			tasks.PUT("/:id", controllers.UpdateTask)
			tasks.DELETE("/:id", controllers.DeleteTask)
		}

	}

}

