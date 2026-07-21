package controllers

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/models"
	"github.com/rechmand/task-management/repository"
)

// ======================================================
// Task Controller
// Berisi implementasi endpoint CRUD Task beserta
// filtering, pagination, Redis Cache dan cache invalidation.
// ======================================================

// Mengambil daftar task berdasarkan parameter filter.
// Data akan diambil dari Redis terlebih dahulu.
// Jika cache tidak ditemukan, data diambil dari MySQL
// kemudian disimpan kembali ke Redis selama 60 detik.
func GetTasks(c *gin.Context) {

	keyword := c.DefaultQuery("keyword", "")
	status := c.DefaultQuery("status", "")
	assignee := c.DefaultQuery("assignee", "")
	sort := c.DefaultQuery("sort", "")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validasi nilai pagination
	// - page minimal 1
	if page < 1 {
		page = 1
	}

	// - limit default 10
	if limit < 1 {
		limit = 10
	}

	// - limit maksimal 100
	if limit > 100 {
		limit = 100
	}

	// Membuat cache key berdasarkan query parameter.
	// Contoh: tasks:status=pending&page=1&limit=10
	ctx := context.Background()
	cacheKey := "tasks:" + c.Request.URL.RawQuery
	cached, err := config.Redis.Get(ctx, cacheKey).Result()

	// Cek apakah data sudah tersedia di Redis.
	// Jika ada, langsung kembalikan response tanpa query database.
	if err == nil {
		c.Data(
			http.StatusOK,
			"application/json",
			[]byte(cached),
		)
		return
	}

	// Cache tidak ditemukan,
	// ambil data terbaru dari database.
	tasks, totalRecords, err := repository.GetAllTasks(keyword, status, assignee, page, limit, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Hitung Total Halaman (Total Pages)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	log.Println("Use Query Database")
	response := gin.H{
		"success":     true,
		"message":     "Tasks retrieved successfully",
		"data":        tasks,
		"total_pages": totalPages,
	}

	// Simpan hasil query ke Redis
	// agar request berikutnya lebih cepat.
	jsonData, _ := json.Marshal(response)
	config.Redis.Set(ctx, cacheKey, jsonData, 60*time.Second)

	// response
	c.JSON(http.StatusOK, response)

}

// Membuat task baru.
// Setelah data berhasil disimpan,
// seluruh cache task akan dihapus
// agar data terbaru selalu ditampilkan.
func CreateTask(c *gin.Context) {

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
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
				"success": false,
				"message": "Task title " + task.Title + " already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return

	}

	// delete redis
	clearTaskCache()

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Task created successfully",
		"data":    task,
	})

}

// get detail
func GetTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	task, err := repository.GetTaskByID(uint(id))

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Task not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task retrieved successfully",
		"data":    task,
	})

}

// Memperbarui data task berdasarkan ID.
// Cache akan dihapus setelah update berhasil.
func UpdateTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	task, err := repository.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Task not found",
		})
		return
	}

	var input models.Task

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
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
				"success": false,
				"message": "Task title " + task.Title + " already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// delete redis
	clearTaskCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task updated successfully",
		"data":    task,
	})
}

// Soft Delete task berdasarkan ID.
// Data tidak benar-benar dihapus dari database,
// melainkan hanya mengisi kolom deleted_at.
func DeleteTask(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	err := repository.DeleteTask(uint(id))
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Task not found",
		})

		return
	}

	// delete redis
	clearTaskCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task deleted successfully",
	})

}

// Menghapus seluruh cache task.
// Dipanggil setiap terjadi Create, Update atau Delete
// agar cache tidak menyimpan data lama.
func clearTaskCache() {
	ctx := context.Background()

	keys, _ := config.Redis.Keys(ctx, "tasks:*").Result()

	for _, key := range keys {
		config.Redis.Del(ctx, key)
	}
}
