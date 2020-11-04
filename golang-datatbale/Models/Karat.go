package Models

import (
	"github.com/jinzhu/gorm"
)

type Karat struct {
	gorm.Model
	Title string	`gorm:"not null" form:"title"`
	Amount float64 `gorm:"not null" form:"amount"`
	Status int	`gorm:"type:tinyint(4); not null; default:0" form:"status"`
	Products []Product	`gorm:"save_associations:false; association_save_reference:false"`
}