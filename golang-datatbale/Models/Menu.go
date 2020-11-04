package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Menu struct {
	gorm.Model
	ParentID uint	`gorm:"not null; default:0"`
	Level int		`gorm:"type:tinyint(4); not null; default:1"`
	Title string	`gorm:"not null"`
	Type int		`gorm:"type:tinyint(4); not null; default:0"`
	Position sql.NullString
	PageUrl	string	`gorm:"default:'#'; not null"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
	Page Page	`gorm:"save_associations:false; association_save_reference:false"`
}