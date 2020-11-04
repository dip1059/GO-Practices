package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type RedsysTransaction struct {
	gorm.Model
	UserID uint	`gorm:"not null"`
	OrderID uint	`gorm:"not null; default:0"`
	RedsysOrder string	`gorm:"index; not null"`
	SubTotal float64	`gorm:"not null; default:0.00"`
	TotalFees float64	`gorm:"not null; default:0.00"`
	TotalVAT float64	`gorm:"not null; default:0.00"`
	Total float64	`gorm:"not null"`
	TotalBids	float64	 `gorm:"not null; default:0.00"`
	PaidStatus int `gorm:"type:tinyint(4); not null; default:0"`
	Refunded sql.NullString
	RedsysOrderDetails []RedsysOrderDetail	`gorm:"save_associations:false; association_save_reference:false"`
	User User	`gorm:"save_associations:false; association_save_reference:false"`
}

