package Models

import (
	"github.com/jinzhu/gorm"
)

type CouponOrder struct {
	gorm.Model
	CouponID uint	`gorm:"not null"`
	UserID uint	    `gorm:"not null"`
	OrderID uint	`gorm:"not null"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
	Coupon Coupon	`gorm:"save_associations:false; association_save_reference:false"`
}