package Models

import (
	"github.com/jinzhu/gorm"
)

type ShippingAddress struct {
	gorm.Model
	UserID     uint
	FirstName  string `gorm:"not null" form:"first_name" binding:"required"`
	LastName   string `gorm:"not null" form:"last_name" binding:"required"`
	Phone      string `gorm:"not null" form:"phone" binding:"required"`
	Country    string `gorm:"not null" form:"country" binding:"required"`
	State   string `gorm:"not null" form:"state"`
	City   string `gorm:"not null" form:"city" binding:"required"`
	ZipCode int `gorm:"not null; default:0" form:"zip_code"`
	Address    string `gorm:"not null" form:"address" binding:"required"`
	Apartment  string `gorm:"not null" form:"apartment"`
	Status int	`gorm:"type:tinyint(4); not null; default:0" form:"status"`
	//Users      User   `gorm:"save_associations:false; association_save_reference:false" json:"-" binding:"-"`
}
