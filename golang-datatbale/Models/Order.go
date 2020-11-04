package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID uint	`gorm:"index; not null"`
	PayMethodID uint	`gorm:"index; not null"`
	BankID sql.NullInt64	`gorm:"type:int(10); index"`
	SubTotal float64	`gorm:"not null; default:0.00"`
	CouponDiscount	float64	 `gorm:"not null; default:0.00"`
	TotalDiscount	float64	 `gorm:"not null; default:0.00"`
	Fees1Percent float64	`gorm:"not null; default:0.00"`
	Fees1Fixed float64	`gorm:"not null; default:0.00"`
	Fees2Percent float64	`gorm:"not null; default:0.00"`
	Fees2Fixed float64	`gorm:"not null; default:0.00"`
	Fees3Percent float64	`gorm:"not null; default:0.00"`
	Fees3Fixed float64	`gorm:"not null; default:0.00"`
	TotalFees float64	`gorm:"not null; default:0.00"`
	GrandTotal float64	`gorm:"not null; default:0.00"`
	TotalGrmAmount	float64	 `gorm:"not null; default:0.00"`
	PaidAmount	float64	 `gorm:"not null; default:0.00"`
	OrderStatus int	`gorm:"type:tinyint(4); not null; default:0"`
	PaymentStatus int	`gorm:"type:tinyint(4); not null; default:0"`
	UniqueCode string	`gorm:"not null; default:'0'"`
	DeliveryType int	`gorm:"type:tinyint(4); not null; default:0"`
	TrackCode string	`gorm:"not null; default:'N/A'"`

	User User	`gorm:"save_associations:false; association_save_reference:false"`
	PayMethod PayMethod	`gorm:"save_associations:false; association_save_reference:false"`
	Bank Bank	`gorm:"save_associations:false; association_save_reference:false"`
	Lusopay Lusopay	`gorm:"save_associations:false; association_save_reference:false"`
	OrderDetails []OrderDetail	`gorm:"save_associations:false; association_save_reference:false"`
	OrderDocs []OrderDoc	`gorm:"save_associations:false; association_save_reference:false"`
	Invoice Invoice	`gorm:"save_associations:false; association_save_reference:false"`
	StripeTransaction StripeTransaction	`gorm:"save_associations:false; association_save_reference:false"`
	CouponOrder CouponOrder	`gorm:"save_associations:false; association_save_reference:false"`
	OrderShippingAddress OrderShippingAddress `gorm:"save_associations:false; association_save_reference:false"`
	MollieOrder MollieOrder	`gorm:"save_associations:false; association_save_reference:false"`
	IsDiscountCoupon bool `gorm:"-"`
	IsFree bool `gorm:"-"`
}


func (order Order) FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func (order Order) FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}