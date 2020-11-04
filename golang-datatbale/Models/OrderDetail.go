package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type OrderDetail struct {
	gorm.Model
	OrderID uint	`gorm:"index; not null"`
	KaratID uint	`gorm:"index; not null"`
	ProductID uint	`gorm:"index; not null"`
	ProductTitle string	`gorm:"not null"`
	ProductType int	`gorm:"not null"`
	ProductKaratAmount float64 `gorm:"not null"`
	ProductGrmAmount float64 `gorm:"not null"`
	TotalGrmAmount	float64	 `gorm:"not null"`
	Quantity float64	`gorm:"not null"`
	ProductPrice float64	`gorm:"not null"`
	ProductImgUrl sql.NullString
	Total float64	`gorm:"not null"`
	TotalDiscount float64	`gorm:"not null"`
	TotalWithDiscount float64	`gorm:"not null"`
	Product Product	`gorm:"save_associations:false; association_save_reference:false"`
}