package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type OrderDoc struct{
	gorm.Model
	OrderID uint	`gorm:"index; not null"`
	DocUrl sql.NullString
	Type int	`gorm:"type:tinyint(4); not null; default:1"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
}