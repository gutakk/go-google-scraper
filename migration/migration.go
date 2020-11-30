package migration

import (
	"fmt"
	"log"

	"github.com/gutakk/go-google-scraper/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(fmt.Sprintf("Failed to migrate %v", err))
	} else {
		log.Print("Migrate user schema successfully")
	}

	if err := db.AutoMigrate(&models.Keyword{}); err != nil {
		log.Fatal(fmt.Sprintf("Failed to migrate %v", err))
	} else {
		log.Print("Migrate keyword schema successfully")
	}
}
