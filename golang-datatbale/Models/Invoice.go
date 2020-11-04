package Models

import (
	"github.com/jinzhu/gorm"
)

type Invoice struct{
	gorm.Model
	OrderID uint	`gorm:"index; not null"`
}
