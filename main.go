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

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.SetupRedisPool()
	oauth.SetupOAuthServer()

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	if error := r.Run(fmt.Sprint(":", os.Getenv("PORT"))); error != nil {
		log.Fatal(fmt.Sprintf("Failed to start the server %v", error))
	}
}
