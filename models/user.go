package models

import "database/sql"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	CreatedAt sql.NullTime `json:"createdAt"`
	UpdateAt  sql.NullTime `json:"updatedAt"`
}
