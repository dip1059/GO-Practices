package Models

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

type WebsiteSetting struct{
	gorm.Model
	ContentName	string	`gorm:"not null; unique_index"`
	Content template.HTML `gorm:"type:text"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}