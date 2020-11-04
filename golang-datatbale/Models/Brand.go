package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Brand struct {
	gorm.Model
	Name string	`gorm:"not null" form:"name"`
	Logo sql.NullString
	Status int	`gorm:"type:tinyint(4); not null; default:0" form:"status"`
	Products []Product	`gorm:"save_associations:false; association_save_reference:false" binding:"-"`
}