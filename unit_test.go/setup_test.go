package unit_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConn    *gorm.DB
	redisConn *redis.Client
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConn, err = gorm.Open(postgres.Open(getPostgresDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	redisConn = redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: getRedisPassword(),
	})

	setupTestData()

	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}

func getPostgresDSN() string {
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable", dbPort, dbUser, dbPassword, dbName)

}

func getRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}

func getRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func setupTestData() {
	cleanup()

	err := dbConn.AutoMigrate(&Article{})
	if err != nil {
		log.Fatalf("Failed to perform database migrations: %v", err)
	}

	article := Article{
		Title:  "Test Article",
		Author: "Elvina Fitriani",
		Body:   "This is a test article.",
	}
	err = dbConn.Create(&article).Error
	if err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}
}

func cleanup() {
	err := dbConn.Exec("TRUNCATE TABLE articles CASCADE").Error
	if err != nil {
		log.Printf("Failed to truncate article table: %v", err)
	}

	err = redisConn.FlushAll(context.Background()).Err()
	if err != nil {
		log.Printf("Failed to clear Redis data: %v", err)
	}
}

type Article struct {
	gorm.Model
	Author string `json:"author" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}
