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

	InitKeywordStatusEnum(db)

	if err := db.AutoMigrate(&models.Keyword{}); err != nil {
		log.Fatal(fmt.Sprintf("Failed to migrate %v", err))
	} else {
		log.Print("Migrate keyword schema successfully")
	}
}

// TODO: Separate file for migration
func InitKeywordStatusEnum(db *gorm.DB) {
	db.Exec(`
		DO $$ BEGIN
			CREATE TYPE keyword_status AS ENUM('pending', 'processing', 'processed', 'failed');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
}
