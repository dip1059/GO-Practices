package Models

import (
	"github.com/jinzhu/gorm"
)

type EmailChangeHistory struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	NewEmail string	`gorm:"not null"`
	OldEmail string	`gorm:"not null"`
	User User	`gorm:"save_associations:false; association_save_reference:false"`
}