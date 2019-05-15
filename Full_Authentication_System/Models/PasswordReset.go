package Models

import "database/sql"

type PasswordReset struct {
	ID int
	Email string
	Token sql.NullString
	Status int
}
