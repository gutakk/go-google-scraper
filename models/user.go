package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique;notNull;index"`
	Password string `gorm:"notNull"`
}
