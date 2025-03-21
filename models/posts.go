package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
}
