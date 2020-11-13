package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	DevMode     = "dev"
	ReleaseMode = "release"
	TestMode    = "test"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV")

	switch env {
	case ReleaseMode:
		err := godotenv.Load(".env." + ReleaseMode)
		if err != nil {
			log.Fatal(fmt.Sprintf("Load .env.release failed with reason %v", err))
		} else {
			log.Print("Load .env.release successfully")
		}
	case TestMode:
		_ = godotenv.Load(".env." + TestMode)
	default:
		_ = godotenv.Load(".env." + DevMode)
	}

	_ = godotenv.Load()
}
