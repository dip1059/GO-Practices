package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type RedsysOrderDetail struct{
	gorm.Model
	RedsysOrder string	`gorm:"index; not null"`
	ProductID uint	`gorm:"index; not null"`
	ProductTitle string	`gorm:"not null" form:"title"`
	ProductGrmAmount float64 `gorm:"not null" form:"grm_amount"`
	Quantity float64	`gorm:"not null"`
	ProductPrice float64	`gorm:"not null" form:"price"`
	ProductImgUrl sql.NullString
	Total float64	`gorm:"not null"`
	TotalTax float64	`gorm:"not null"`
	TotalWithTax float64	`gorm:"not null"`
	//Product Product	`gorm:"save_associations:false; association_save_reference:false"`
}