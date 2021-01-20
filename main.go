package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"
)

func main() {
	config.LoadEnv()
	database := db.ConnectDB()
	migration.Migrate(database)

	db.SetupRedisPool()

	r := config.SetupRouter()
	controllers.CombineRoutes(r)

	if error := r.Run(fmt.Sprint(":", os.Getenv("PORT"))); error != nil {
		glog.Fatalf("Failed to start the server %s", error)
	}
}
