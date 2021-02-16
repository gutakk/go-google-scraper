package main

import (
	"fmt"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/migration"
	"github.com/gutakk/go-google-scraper/oauth"
)

const (
	startOAuthServerFailureError = "Failed to start oauth server: "
	startServerFailureError      = "Failed to start the server: "
)

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.SetupRedisPool()
	err := oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(startOAuthServerFailureError, err)
	}

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	err = r.Run(fmt.Sprint(":", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(startServerFailureError, err)
	}
}
