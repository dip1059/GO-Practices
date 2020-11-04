package Models

import (
	"github.com/jinzhu/gorm"
	"html/template"
)

type HomeContent struct {
	gorm.Model
	TextContent template.HTML	`gorm:"type:longtext"`
	Image string	`gorm:"not null"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
}