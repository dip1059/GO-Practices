package Models

import (
	"github.com/jinzhu/gorm"
)

type ProcessedTransfer struct {
	gorm.Model
	TransferID string `gorm:"not null; unique"`
}
