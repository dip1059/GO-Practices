package Models

import "database/sql"

type User struct {
	ID int
	Name string
	Email string
	Password string
	Role int
	ActiveStauts int
	EmailVf sql.NullString
}

