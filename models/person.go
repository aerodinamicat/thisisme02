package models

import (
	"time"
)

type Person struct {
	FirstName     string    `json:"firstName"`
	SecondName    string    `json:"secondName"`
	FirstSurname  string    `json:"firstSurname"`
	SecondSurname string    `json:"secondSurname"`
	Gender        string    `json:"gender"`
	BirthDate     time.Time `json:"birthDate"`

	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
