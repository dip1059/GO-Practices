package Models

import (
	"github.com/jinzhu/gorm"
)

type Newsletter struct{
	gorm.Model
	Email string	`gorm:"not null; unique_index" form:"email" binding:"email"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
}