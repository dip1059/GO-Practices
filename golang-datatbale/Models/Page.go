package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"html/template"
)

type Page struct{
	gorm.Model
	Url	string	`gorm:"not null; unique_index"`
	Title string	`gorm:"not null; unique_index"`
	TextContent template.HTML `gorm:"type:longtext"`
	ImgUrl sql.NullString
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}