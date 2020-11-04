package Models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Coupon struct {
	gorm.Model
	Title string `gorm:"not null"`
	Code string	`gorm:"not null"`
	Type int	`gorm:"type:tinyint(4); not null"`
	Validity float32 `gorm:"default:0.0"`
	Amount float64 `gorm:"not null"`
	StartDate time.Time
	EndDate time.Time
	Description string `gorm:"type:text"`
	Status int	`gorm:"type:tinyint(4); not null; default:0"`
}