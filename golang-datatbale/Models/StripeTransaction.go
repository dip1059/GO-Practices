package Models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

type StripeTransaction struct {
	gorm.Model
	OrderID uint	`gorm:"index; not null"`
	StripeToken string	`gorm:"not null; unique_index"`
	ChargeID string		`gorm:"not null; unique_index"`
	CustomerEmail string   `gorm:"not null"`
	Currency string	`gorm:"not null"`
	Amount int	`gorm:"not null"`
	PaidStatus int `gorm:"type:tinyint(4); not null; default:0"`
	Status	string `gorm:"not null; default:'pending'"`
	Refunded sql.NullString
	//Order Order	`gorm:"save_associations:false; association_save_reference:false"`
}

