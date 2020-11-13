package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	r := config.SetupRouter()

	if error := r.Run(":" + os.Getenv("PORT")); error != nil {
		log.Fatal(fmt.Sprintf("Failed to start the server %v", error))
	}
}
