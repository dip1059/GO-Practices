package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type AdminSetting struct{
	gorm.Model
	Slug string 	`gorm:"unique_index; not null" form:"slug"`
	Value sql.NullString	`gorm:"type:text"`
}