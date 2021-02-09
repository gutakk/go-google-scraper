package main

import (
	"fmt"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"
	"github.com/gutakk/go-google-scraper/oauth"

	log "github.com/sirupsen/logrus"
)

const (
	startOAuthServerFailureError = "Failed to start oauth server: %v"
	startServerFailureError      = "Failed to start the server: %v"
)

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.SetupRedisPool()
	err := oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(fmt.Sprintf(startOAuthServerFailureError, err))
	}

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	err = r.Run(fmt.Sprint(":", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(fmt.Sprintf(startServerFailureError, err))
	}
}
