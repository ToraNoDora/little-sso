package models

import "time"

type Group struct {
	ID          int    `json:"-" db:"id"`
	AppID       string `json:"-" db:"app_id"`
	Name        string `json:"-" db:"name"`
	Description string `json:"-" db:"description"`
}

type GroupRole struct {
	ID        int       `json:"-" db:"id"`
	GroupID   int       `json:"-" db:"group_id"`
	RoleID    int       `json:"-" db:"role_id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}
