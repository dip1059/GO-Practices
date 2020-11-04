package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Title string	`gorm:"not null" form:"title"`
	KaratID uint `gorm:"not null" form:"karat_id"`
	BrandID uint `gorm:"not null" form:"brand_id"`
	GrmAmount float64 `gorm:"not null" form:"grm_amount"`
	Type int `gorm:"type:tinyint(4); not null; default:1" form:"type"`
	Min string `gorm:"type:varchar(40); not null; default:'1'" form:"min"`
	Max string `gorm:"type:varchar(40); not null; default:'9999999999999'" form:"max"`
	Description sql.NullString	`gorm:"type:text"`
	Price float64	`gorm:"not null" form:"price"`
	Discount float64	`gorm:"not null; default:0.0" form:"discount"`
	Status int	`gorm:"type:tinyint(4); not null; default:0" form:"status"`
	ImgUrl sql.NullString
	Position int 	`gorm:"type:tinyint(4); not null; default:1" form:"position"`
	//SecondTitle string	`gorm:"not null; default:''" form:"second_title"`
	Karat Karat	`gorm:"save_associations:false; association_save_reference:false"`
	Brand Brand	`gorm:"save_associations:false; association_save_reference:false"`
	DiscountPrice float64	`gorm:"-"`
	IsInWishlist bool	`gorm:"-"`
}