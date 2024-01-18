package models

type Role struct {
	ID          int    `json:"-" db:"id"`
	Name        string `json:"-" db:"name"`
	Description string `json:"-" db:"description"`
}
