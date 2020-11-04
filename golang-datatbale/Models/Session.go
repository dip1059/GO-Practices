package Models

import "database/sql"

type Session struct {
	ID string	`gorm:"primary_key"`
	UserID sql.NullInt64
	IpAddress sql.NullString
	UserAgent sql.NullString
	Payload string
	LastActivity int
	User User
}
