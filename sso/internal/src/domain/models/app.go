package models

import "time"

type App struct {
	ID          string    `json:"-" db:"id"`
	Description string    `json:"description" db:"description"`
	Name        string    `json:"name" db:"name"`
	Secret      string    `json:"secret" db:"secret"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
