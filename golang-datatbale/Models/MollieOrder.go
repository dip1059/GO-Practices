package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type MollieOrder struct {
	gorm.Model
	UserID uint	`gorm:"not null"`
	OrderID uint	`gorm:"not null; default:0"`
	PaymentID string	`gorm:"index; not null"`
	SubTotal float64	`gorm:"not null; default:0.00"`
	CouponDiscount	float64	 `gorm:"not null; default:0.00"`
	TotalDiscount	float64	 `gorm:"not null; default:0.00"`
	Fees float64	`gorm:"not null; default:0.00"`
	VAT float64	`gorm:"not null; default:0.00"`
	TotalFees float64	`gorm:"not null; default:0.00"`
	TotalVAT float64	`gorm:"not null; default:0.00"`
	GrandTotal float64	`gorm:"not null; default:0.00"`
	TotalGrmAmount	float64	 `gorm:"not null; default:0.00"`
	PaymentStatus string `gorm:"not null; default:''"`
	Refunded sql.NullString
	PaymentUrl string `gorm:"not null; default:''"`
	PaymentMethod string `gorm:"not null; default:''"`
	RedsysOrderDetails []RedsysOrderDetail	`gorm:"save_associations:false; association_save_reference:false"`
	User User	`gorm:"save_associations:false; association_save_reference:false"`
}

