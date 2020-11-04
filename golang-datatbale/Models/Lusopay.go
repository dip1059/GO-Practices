package Models

import (
	"github.com/jinzhu/gorm"
)

type Lusopay struct{
	gorm.Model
	OrderID uint `gorm:"not null"`
	Entity string	`gorm:"not null"`
	Amount float64	`gorm:"not null"`
	Reference string `gorm:"not null"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}

