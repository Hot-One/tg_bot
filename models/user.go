package models

import (
	"time"
)

type User struct {
	ID             string `json:"id" db:"id"`
	FirstName      string `json:"first_name" db:"first_name"`
	LastName       string `json:"last_name" db:"last_name"`
	PhoneNumber    string `json:"phone_number" db:"phone_number"`
	TelegramUserID int    `json:"telegram_user_id" db:"telegram_user_id"`
	Status         Status `json:"status" db:"status"`
	UserName       string `json:"user_name" db:"user_name"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Status string

const (
	StatusUnverified Status = "unverified"
	StatusActive     Status = "active"
	StatusBannded    Status = "banned"

	StatusNew       Status = "new"
	StatusInProcess Status = "in_process"
	StatusDone      Status = "done"
)
