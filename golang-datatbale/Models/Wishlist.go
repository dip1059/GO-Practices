package Models

import (
	"github.com/jinzhu/gorm"
)

type Wishlist struct{
	gorm.Model
	UserID uint	`gorm:"index; not null"`
	ProductID uint	`gorm:"index; not null"`
	Product Product	`gorm:"save_associations:false; association_save_reference:false"`
}