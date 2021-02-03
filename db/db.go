package db

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (db *gorm.DB) {
	db, err := gorm.Open(postgres.Open(constructDsn()), &gorm.Config{})

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database %v", err))
	} else {
		log.Print("Connect to database successfully")
	}

	DB = db

	return db
}

func constructDsn() string {
	if gin.Mode() == gin.ReleaseMode {
		return os.Getenv("DATABASE_URL")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s password=%s",
		host,
		port,
		dbName,
		username,
		password,
	)
}

func GetDatabaseURL() string {
	if gin.Mode() == gin.ReleaseMode {
		return os.Getenv("DATABASE_URL")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
}

var GetDB = func() *gorm.DB {
	return DB
}
