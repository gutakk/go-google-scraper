package migration

import (
	"github.com/golang/glog"
	"github.com/gutakk/go-google-scraper/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		glog.Fatalf("Failed to migrate %s", err)
	} else {
		glog.Info("Migrate user schema successfully")
	}

	InitKeywordStatusEnum(db)

	if err := db.AutoMigrate(&models.Keyword{}); err != nil {
		glog.Fatalf("Failed to migrate %s", err)
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
