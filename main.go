package main

import (
	"fmt"
	"log"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/migration"
)

func main() {
	config.LoadEnv()
	db := db.ConnectDB()
	migration.Migrate(db)

	r := config.SetupRouter()

	if error := r.Run(); error != nil {
		log.Fatal(fmt.Sprintf("Failed to start the server %v", error))
	}
}
