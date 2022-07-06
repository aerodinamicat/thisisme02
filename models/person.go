package models

import "database/sql"

type Person struct {
	FirstName     string       `json:"firstName"`
	SecondName    string       `json:"secondName"`
	FirstSurname  string       `json:"firstSurname"`
	SecondSurname string       `json:"secondSurname"`
	Gender        string       `json:"gender"`
	BirthDate     sql.NullTime `json:"birthDate"`

	UserId    string       `json:"userId"`
	CreatedAt sql.NullTime `json:"createdAt"`
	UpdatedAt sql.NullTime `json:"updatedAt"`
}
