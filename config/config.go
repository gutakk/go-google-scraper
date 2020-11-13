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

	if "" == env {
		env = DevMode
	}

	godotenv.Load(".env." + env + ".local")
	if TestMode != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}
