package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/models"
	"github.com/rechmand/task-management/repository"
)

// get task
func GetTasks(c *gin.Context) {

	keyword := c.DefaultQuery("keyword", "")
	status := c.DefaultQuery("status", "")
	assignee := c.DefaultQuery("assignee", "")
	sort := c.DefaultQuery("sort", "")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// validasi page
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	// ========== redis open
	ctx := context.Background()
	cacheKey := "tasks:" + c.Request.URL.RawQuery

	// =========================
	// Check Redis
	// =========================
	cached, err := config.Redis.Get(ctx, cacheKey).Result()

	if err == nil {

		c.Data(
			http.StatusOK,
			"application/json",
			[]byte(cached),
		)

		return
	}
	// ========== redis close

	// get data from database
	tasks, err := repository.GetAllTasks(keyword, status, assignee, page, limit, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println("Use Query Database")
	response := gin.H{
		"success": true,
		"data":    tasks,
	}

	// =========================
	// Save Redis
	// =========================
	jsonData, _ := json.Marshal(response)
	config.Redis.Set(ctx, cacheKey, jsonData, 60*time.Second)

	// response
	c.JSON(http.StatusOK, response)

}

// post task
func CreateTask(c *gin.Context) {

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// create data
	err := repository.CreateTask(&task)
	if err != nil {

		// jika error duplicate title
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Task title " + task.Title + " already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return

	}

	// delete redis
	ctx := context.Background()
	keys, _ := config.Redis.Keys(ctx, "tasks:*").Result()
	for _, key := range keys {
		config.Redis.Del(ctx, key)
	}

	c.JSON(http.StatusCreated, task)

}

// get detail
func GetTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	task, err := repository.GetTaskByID(uint(id))

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})

		return
	}

	c.JSON(http.StatusOK, task)

}

// update task
func UpdateTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	task, err := repository.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})
		return
	}

	var input models.Task

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Assignee = input.Assignee

	if err := repository.UpdateTask(&task); err != nil {

		// jika error duplicate title
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Task title " + task.Title + " already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// delete redis
	ctx := context.Background()
	keys, _ := config.Redis.Keys(ctx, "tasks:*").Result()
	for _, key := range keys {
		config.Redis.Del(ctx, key)
	}

	c.JSON(http.StatusOK, task)
}

// delete
func DeleteTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	err := repository.DeleteTask(uint(id))

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
		})

		return
	}

	// delete redis
	ctx := context.Background()
	keys, _ := config.Redis.Keys(ctx, "tasks:*").Result()
	for _, key := range keys {
		config.Redis.Del(ctx, key)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})

}
