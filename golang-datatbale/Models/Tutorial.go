package Models

import (
	"github.com/jinzhu/gorm"
)

type Tutorial struct {
	gorm.Model
	TutorialType int `gorm:"not null;"`
	Language string	`gorm:"not null;"`
	Link string	`gorm:"not null;"`
}
