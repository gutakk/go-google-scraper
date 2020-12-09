package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"
)

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.GenerateRedisPool("localhost:6379")

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	if error := r.Run(fmt.Sprint(":", os.Getenv("APP_PORT"))); error != nil {
		log.Fatal(fmt.Sprintf("Failed to start the server %v", error))
	}
}
