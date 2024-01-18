package models

import "time"

type Admin struct {
	ID        string    `json:"-" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	AppID     string    `json:"app_id" db:"app_id"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
