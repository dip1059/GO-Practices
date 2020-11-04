package Models

import (
	"github.com/jinzhu/gorm"
)

type WalletHistory struct {
	gorm.Model
	WalletID uint	`gorm:"not null"`
	GoldAmountBefore float64 `gorm:"not null; default:0.0"`
	AddedGoldAmount float64 `gorm:"not null; default:0.0"`
	GoldAmountAfter float64	`gorm:"not null; default:0.0"`
	Note string
	Status int	`gorm:"type:tinyint(4); not null; default:1"`
}