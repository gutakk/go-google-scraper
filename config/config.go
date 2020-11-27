package config

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	switch gin.Mode() {
	case gin.ReleaseMode:
		_ = godotenv.Load(".env." + gin.ReleaseMode)
	case gin.TestMode:
		_ = godotenv.Load(".env." + gin.TestMode)
	default:
		_ = godotenv.Load(".env." + gin.DebugMode)
	}

	_ = godotenv.Load()
}
