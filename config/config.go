package config

import (
	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	switch gin.Mode() {
	case gin.ReleaseMode:
		err := godotenv.Load(".env." + gin.ReleaseMode)
		if err != nil {
			log.Error("Failed to load release env: ", err)
		}
	case gin.TestMode:
		err := godotenv.Load(".env." + gin.TestMode)
		if err != nil {
			log.Error("Failed to load test env: ", err)
		}
	default:
		err := godotenv.Load(".env." + gin.DebugMode)
		if err != nil {
			log.Error("Failed to load debug env: ", err)
		}
	}

	err := godotenv.Load()
	if err != nil {
		log.Error("Failed to load env: ", err)
	}
}
