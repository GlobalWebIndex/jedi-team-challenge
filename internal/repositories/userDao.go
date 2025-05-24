package repositories

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;"`
	Password string `gorm:"not null;"`
}
