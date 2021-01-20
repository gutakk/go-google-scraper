package migration

import (
	"github.com/gutakk/go-google-scraper/models"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	migrateUserErr := db.AutoMigrate(&models.User{})
	if migrateUserErr != nil {
		glog.Fatalf("Failed to migrate %s", migrateUserErr)
	} else {
		glog.Info("Migrate user schema successfully")
	}

	InitKeywordStatusEnum(db)

	migrateKeywordErr := db.AutoMigrate(&models.Keyword{})
	if migrateKeywordErr != nil {
		glog.Fatalf("Failed to migrate %s", migrateKeywordErr)
	} else {
		glog.Info("Migrate keyword schema successfully")
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
