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
		godotenv.Load(".env." + ReleaseMode)
	case TestMode:
		godotenv.Load(".env." + TestMode)
	default:
		godotenv.Load(".env." + DevMode)
	}

	_ = godotenv.Load()
}
