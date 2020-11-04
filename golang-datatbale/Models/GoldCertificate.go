package Models

import (
	"time"
)

type GoldCertificate struct {
	ID        string `gorm:"primary_key; type:varchar(36)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
	UserID uint `gorm:"not null"`
	OrderID uint `gorm:"not null"`
	DeliveryType int	`gorm:"type:tinyint(4); not null; default:1"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
}