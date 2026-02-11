package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestSuccess(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		data := map[string]string{"key": "value"}
		Success(c, data)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", response.Message)
	assert.Equal(t, map[string]interface{}{"key": "value"}, response.Data)
}

func TestCreated(t *testing.T) {
	r := setupTestRouter()
	r.POST("/test", func(c *gin.Context) {
		data := map[string]string{"id": "1"}
		Created(c, data)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, "created successfully", response.Message)
	assert.Equal(t, map[string]interface{}{"id": "1"}, response.Data)
}

func TestBadRequest(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		BadRequest(c, "Invalid request")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "Invalid request", response.Message)
	assert.Equal(t, "Invalid request", response.Error)
}

func TestNotFound(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		NotFound(c, "Resource not found")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, "Resource not found", response.Message)
	assert.Equal(t, "Resource not found", response.Error)
}

func TestInternalServerError(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		InternalServerError(c, "Database error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Equal(t, "Database error", response.Message)
	assert.Equal(t, "Database error", response.Error)
}

func TestValidationError(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		ValidationError(c, "Invalid email format")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "validation error: Invalid email format", response.Message)
	assert.Equal(t, "validation error: Invalid email format", response.Error)
}

func TestDatabaseError(t *testing.T) {
	r := setupTestRouter()
	r.GET("/test", func(c *gin.Context) {
		DatabaseError(c, "Connection failed")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 500, response.Code)
	assert.Equal(t, "database error: Connection failed", response.Message)
	assert.Equal(t, "database error: Connection failed", response.Error)
}
