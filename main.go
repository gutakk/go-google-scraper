package main

import (
	"fmt"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"

	"github.com/golang/glog"
)

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.SetupRedisPool()

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	err := r.Run(fmt.Sprint(":", os.Getenv("PORT")))
	if err != nil {
		glog.Fatalf("Failed to start the server %s", err)
	}
}
