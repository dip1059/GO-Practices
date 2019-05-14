package Models

import (
	"database/sql"
)

type User struct {
	ID int								`field_name:"id"`
	Name string							`field_name:"full_name"`
	Email string						`field_name:"email"`
	Phone sql.NullString				`field_name:"phone"`
	Password string						`field_name:"password"`
	ActiveStauts int					`field_name:"activestatus"`
	Role int							`field_name:"role_id"`
	EmailVf sql.NullString				`field_name:"email_verification"`
	RememberToken sql.NullString		`field_name:"remember_token"`
	Deleted_At sql.NullString			`field_name:"deleted_at"`
	Created_At sql.NullString			`field_name:"created_at"`
	Updated_At sql.NullString			`field_name:"deleted_at"`
}

