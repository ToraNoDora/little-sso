package models

import "time"

type User struct {
	ID          string    `json:"-" db:"id"`
	Username    string    `json:"-" db:"username"`
	Email       string    `json:"-" db:"email"`
	PassHash    []byte    `json:"-" db:"pass_hash"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	Permissions []Permission
}

type Permission struct {
	ID        string    `json:"-" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	GroupID   int       `json:"group_id" db:"group_id"`
	AddFlag   bool      `json:"add_flag" db:"add_flag"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}
