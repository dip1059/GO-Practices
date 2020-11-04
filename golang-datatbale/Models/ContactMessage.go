package Models

import (
	"github.com/jinzhu/gorm"
)

type ContactMessage struct{
	gorm.Model
	FullName string	`gorm:"not null" form:"full_name"`
	Email string	`gorm:"not null" form:"email"`
	Phone string	`gorm:"not null" form:"phone"`
	Message string	`gorm:"type:longtext; not null" form:"message"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}