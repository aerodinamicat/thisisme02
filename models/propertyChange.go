package models

import "time"

type PropertyChange struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	From   string `json:"from"`
	To     string `json:"to"`

	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}
