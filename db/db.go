package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"huinya/models"
	"log"
	"os"
)

var DB *gorm.DB

func DatabaseConfiguration() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=disable password=%s",
		dbHost, dbUser, dbName, dbPort, dbPassword,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	if DB == nil {
		log.Fatal("Database connection failed")
	}

	DB.AutoMigrate(&models.User{}, &models.Post{})
}
