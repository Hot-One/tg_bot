package models

import (
	"time"
)

type Order struct {
	ID      string `json:"id" db:"id"`
	UserID  string `json:"user_id" db:"user_id"`
	Lat     string `json:"lat" db:"lat"`
	Long    string `json:"long" db:"long"`
	Block   string `json:"building" db:"building"`
	Floor   string `json:"floor" db:"floor"`
	Flat    string `json:"flat" db:"flat"`
	Photo   string `json:"photo" db:"photo"`
	Weight  string `json:"weight" db:"weight"`
	Comment string `json:"comment" db:"comment"`
	Status  Status `json:"status" db:"status"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
