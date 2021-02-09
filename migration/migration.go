package migration

import (
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(errorHandler.MigrateDatabaseFailure, err)
	} else {
		log.Print("Migrate user schema successfully")
	}

	InitKeywordStatusEnum(db)

	err = db.AutoMigrate(&models.Keyword{})
	if err != nil {
		log.Fatal(errorHandler.MigrateDatabaseFailure, err)
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
