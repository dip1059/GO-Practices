package Models

import (
	"github.com/jinzhu/gorm"
)

type PasswordReset struct {
	gorm.Model
	Email string	`gorm:"index; not null"`
	Token string
	Code string
	Status int	`gorm:"type:tinyint(4); not null"`
	NewPassword string `gorm:"-"`
}

