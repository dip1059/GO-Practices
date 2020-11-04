package Models

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

type Bank struct{
	gorm.Model
	Name string	`gorm:"not null"`
	Details template.HTML `gorm:"type:text"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}