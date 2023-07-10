package unit_test

import (
	"article/routers"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupSuccess(t *testing.T) {
	if dbConn == nil {
		t.Errorf("Failed to connect to PostgreSQL database")
	}

	_, err := redisConn.Ping(context.Background()).Result()
	if err != nil {
		t.Errorf("Failed to connect to Redis server: %v", err)
	}

	var count int64
	err = dbConn.Model(&Article{}).Count(&count).Error
	if err != nil {
		t.Errorf("Failed to retrieve article count: %v", err)
	}
	if count != 1 {
		t.Errorf("Invalid article count. Expected: 1, Got: %d", count)
	}
}

func TestCreate(t *testing.T) {
	db := dbConn
	redisClient := redisConn

	article := Article{
		Author: "Test Author",
		Title:  "Test Title",
		Body:   "This is a test body article.",
	}

	router := gin.Default()
	routes := routers.Routes{
		Db:    db,
		R:     router,
		Redis: redisClient,
	}
	routes.Routers()

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(article)
	req, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var savedArticle Article
	db.Last(&savedArticle)
	assert.Equal(t, article.Title, savedArticle.Title)
	assert.Equal(t, article.Author, savedArticle.Author)
	assert.Equal(t, article.Body, savedArticle.Body)

	cacheKey := fmt.Sprintf("article:%d", savedArticle.ID)

	jsonVal, _ := json.Marshal(savedArticle)
	err := redisClient.Set(context.Background(), cacheKey, jsonVal, 0).Err()
	assert.NoError(t, err)

	redisData, err := redisClient.Get(context.Background(), cacheKey).Result()
	assert.NoError(t, err)
	var redisArticle Article
	json.Unmarshal([]byte(redisData), &redisArticle)

	assert.Equal(t, savedArticle.Title, redisArticle.Title)
	assert.Equal(t, savedArticle.Author, redisArticle.Author)
	assert.Equal(t, savedArticle.Body, redisArticle.Body)
}

func TestGetAll(t *testing.T) {
	db := dbConn
	redisClient := redisConn

	router := gin.Default()
	routes := routers.Routes{
		Db:    db,
		R:     router,
		Redis: redisClient,
	}
	routes.Routers()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/article", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Articles []Article `json:"Articles"`
		Response string    `json:"Response"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.Articles)
	assert.Equal(t, "Data retrieved successfully.", response.Response)

}

func TestGetByKeyword(t *testing.T) {
	db := dbConn
	redisClient := redisConn

	router := gin.Default()
	routes := routers.Routes{
		Db:    db,
		R:     router,
		Redis: redisClient,
	}
	routes.Routers()

	article := Article{
		Title:  "Test Article",
		Author: "Elvina Fitriani",
		Body:   "This is a test article.",
	}

	reqBody, _ := json.Marshal(article)
	req, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", fmt.Sprintf("/article/search/keyword/%s", article.Title), nil)

	router.ServeHTTP(w, reqGet)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Articles []Article `json:"Articles"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response.Articles)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Articles))
	assert.Equal(t, article.Title, response.Articles[0].Title)
	assert.Equal(t, article.Author, response.Articles[0].Author)
	assert.Equal(t, article.Body, response.Articles[0].Body)

}

func TestGetByAuthor(t *testing.T) {
	db := dbConn
	redisClient := redisConn

	router := gin.Default()
	routes := routers.Routes{
		Db:    db,
		R:     router,
		Redis: redisClient,
	}
	routes.Routers()

	article := Article{
		Title:  "Test Article",
		Author: "Elvina Fitriani",
		Body:   "This is a test article.",
	}

	reqBody, _ := json.Marshal(article)
	req, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", fmt.Sprintf("/article/search/author/%s", article.Author), nil)

	router.ServeHTTP(w, reqGet)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Articles []Article `json:"Articles"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response.Articles)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Articles))
	assert.Equal(t, article.Title, response.Articles[0].Title)
	assert.Equal(t, article.Author, response.Articles[0].Author)
	assert.Equal(t, article.Body, response.Articles[0].Body)

}
