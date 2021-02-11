package tests

import (
	"fmt"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDatabase() *gorm.DB {
	database, err := gorm.Open(postgres.Open(constructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorconf.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	db.SetupRedisPool()

	initKeywordStatusEnum(db.GetDB())

	err = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if err != nil {
		log.Fatal(errorconf.MigrateDatabaseFailure, err)
	}

	return db.GetDB()
}

func constructTestDsn() string {
	host := "localhost"
	port := "5433"
	dbName := "go_google_scraper_test"
	username := "postgres"

	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s",
		host,
		port,
		dbName,
		username,
	)
}

func initKeywordStatusEnum(db *gorm.DB) {
	db.Exec(`
		DO $$ BEGIN
			CREATE TYPE keyword_status AS ENUM('pending', 'processing', 'processed', 'failed');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
}

func RedisKeyJobs(namespace, jobName string) string {
	return redisKeyJobsPrefix(namespace) + jobName
}

func redisKeyJobsPrefix(namespace string) string {
	return redisNamespacePrefix(namespace) + "jobs:"
}

func redisNamespacePrefix(namespace string) string {
	l := len(namespace)
	if (l > 0) && (namespace[l-1] != ':') {
		namespace = namespace + ":"
	}
	return namespace
}
