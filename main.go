package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"
	"github.com/gutakk/go-google-scraper/oauth"
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
	oauthServerErr := oauth.SetupOAuthServer()
	if oauthServerErr != nil {
		log.Fatal(fmt.Sprintf(startOAuthServerFailureError, oauthServerErr))
	}

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	if error := r.Run(fmt.Sprint(":", os.Getenv("PORT"))); error != nil {
		log.Fatal(fmt.Sprintf(startServerFailureError, error))
	}
}
