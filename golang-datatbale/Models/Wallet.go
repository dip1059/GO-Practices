package Models

import (
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	gorm.Model
	UserID uint	`gorm:"not null"`
	GoldAmount float64	`gorm:"type:double unsigned; not null; default:0.0"`
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
	Users []User  `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
	WalletHistories []WalletHistory  `gorm:"save_associations:false; association_save_reference:false" binding:"-" json:"-"`
}