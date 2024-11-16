package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mubashshir/go-crud/controllers"
	"github.com/mubashshir/go-crud/initializers"
	"github.com/mubashshir/go-crud/models"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/posts", controllers.PostsCreate)
	r.GET("/posts", controllers.PostsIndex)
	r.GET("/posts/:id", controllers.PostsShow)
	r.PUT("/posts/:id", controllers.PostsUpdate)
	r.DELETE("/posts/:id", controllers.PostsDelete)
	return r
}

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.DB.AutoMigrate(&models.Post{})
}

func resetDatabase() {
	initializers.DB.Exec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE")
}

func TestPostsCRUD(t *testing.T) {
	resetDatabase()

	router := setupRouter()

	t.Run("Create Post", func(t *testing.T) {
		body := map[string]string{
			"title": "Test Title",
			"body":  "Test Body",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, body["title"], response["post"].(map[string]interface{})["title"])
		assert.Equal(t, body["body"], response["post"].(map[string]interface{})["body"])
	})

	t.Run("Get All Posts", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/posts", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Greater(t, len(response["posts"].([]interface{})), 0)
	})

	t.Run("Get Single Post", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/posts/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.NotNil(t, response["post"])
	})

	t.Run("Update Post", func(t *testing.T) {
		body := map[string]string{
			"title": "Updated Title",
			"body":  "Updated Title",
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/posts/1", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, body["title"], response["post"].(map[string]interface{})["title"])
		assert.Equal(t, body["body"], response["post"].(map[string]interface{})["body"])
	})

	t.Run("Delete Post", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/posts/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Code)

		req, _ = http.NewRequest("GET", "/posts/1", nil)
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, 404, resp.Code)
	})
}
