package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type BrainTreeTransaction struct {
	gorm.Model
	OrderID uint	`gorm:"index; not null"`
	TxID string		`gorm:"not null; unique_index"`
	CustomerEmail string   `gorm:"not null"`
	PaymentInstrumentType string	`gorm:"not null"`
	Currency string	`gorm:"not null"`
	Amount int	`gorm:"not null"`
	Status	string `gorm:"not null; default:'pending'"`
	Refunded sql.NullString
	Order Order	`gorm:"save_associations:false; association_save_reference:false"`
}

