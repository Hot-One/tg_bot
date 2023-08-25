package models

import "time"

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
	Status  string `json:"status" db:"status"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type User struct {
	ID             string `json:"id" db:"id"`
	FirstName      string `json:"first_name" db:"first_name"`
	LastName       string `json:"last_name" db:"last_name"`
	PhoneNumber    string `json:"phone_number" db:"phone_number"`
	TelegramUserID string `json:"telegram_user_id" db:"telegram_user_id"`
	Status         string `json:"status" db:"status"`
	UserName       string `json:"user_name" db:"user_name"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
