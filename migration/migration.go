package migration

import (
	"github.com/gutakk/go-google-scraper/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	migrateUserErr := db.AutoMigrate(&models.User{})
	if migrateUserErr != nil {
		log.Fatalf("Failed to migrate %s", migrateUserErr)
	} else {
		log.Info("Migrate user schema successfully")
	}

	InitKeywordStatusEnum(db)

	migrateKeywordErr := db.AutoMigrate(&models.Keyword{})
	if migrateKeywordErr != nil {
		log.Fatalf("Failed to migrate %s", migrateKeywordErr)
	} else {
		log.Info("Migrate keyword schema successfully")
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
