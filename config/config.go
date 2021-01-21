package config

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func LoadEnv() {
	var err error
	switch gin.Mode() {
	case gin.ReleaseMode:
		err = godotenv.Load(".env." + gin.ReleaseMode)
	case gin.TestMode:
		err = godotenv.Load(".env." + gin.TestMode)
	default:
		err = godotenv.Load(".env." + gin.DebugMode)
	}

	if err != nil {
		log.Errorf("Load %s env error: %s", gin.Mode(), err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Errorf("Load env error: %s", err)
	}
}
