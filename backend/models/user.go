package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;index:idx_user_username"`
	Password string `gorm:"not null"`
	Level    int    `gorm:"not null;index:idx_user_level"`
	Avatar   string
}
