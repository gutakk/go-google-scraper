package config

import (
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
		_ = godotenv.Load(".env." + ReleaseMode)
	case TestMode:
		_ = godotenv.Load(".env." + TestMode)
	default:
		_ = godotenv.Load(".env." + DevMode)
	}

	_ = godotenv.Load()
}
