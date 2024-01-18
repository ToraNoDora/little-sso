package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, username string, email string, passHash []byte) (string, error) {
	const op = "storage.postgres.CreateUser"

	query := fmt.Sprintf("INSERT INTO %s (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRowContext(ctx, query, username, email, passHash)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
