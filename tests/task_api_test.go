package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rechmand/task-management/config"
	"github.com/rechmand/task-management/routes"
	"github.com/stretchr/testify/assert"
)

// setup initializes the environment, database, and Redis connection used by tests.
func setup() {
	config.LoadEnv()
	config.ConnectDatabase()
	config.ConnectRedis()
}

// TestSearchTask checks that the task search endpoint returns HTTP 200 OK
// when a keyword query is provided.
func TestSearchTask(t *testing.T) {

	setup()

	// Create router and register routes for the test.
	router := gin.Default()

	routes.SetupRoutes(router)

	// Build the GET request for task search.
	req, err := http.NewRequest(
		"GET",
		"/api/tasks?keyword=test",
		nil,
	)

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUpdateTask checks that updating a task returns HTTP 200 OK.
func TestUpdateTask(t *testing.T) {

	setup()

	// Create router and register routes for the test.
	router := gin.Default()
	routes.SetupRoutes(router)

	// JSON payload for updating task with id 1.
	body := `{
		"title":"Updated Task",
		"description":"Updated Description",
		"status":"done",
		"assignee":"Rechmand"
	}`

	req, _ := http.NewRequest(
		"PUT",
		"/api/tasks/1",
		bytes.NewBuffer([]byte(body)),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRedisCache verifies the tasks endpoint works on repeated requests.
func TestRedisCache(t *testing.T) {

	setup()

	// Create router and register routes for the test.
	router := gin.Default()
	routes.SetupRoutes(router)

	// First request to retrieve all tasks.
	req, _ := http.NewRequest(
		"GET",
		"/api/tasks",
		nil,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	assert.Equal(t, http.StatusOK, w2.Code)
}
