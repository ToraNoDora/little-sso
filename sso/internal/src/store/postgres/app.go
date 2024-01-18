package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
)

type AppPostgres struct {
	db *sqlx.DB
}

func NewAppPostgres(db *sqlx.DB) *AppPostgres {
	return &AppPostgres{db: db}
}

func (a *AppPostgres) GetApp(ctx context.Context, appID string) (models.App, error) {
	const op = "storage.postgres.App"

	query := fmt.Sprintf("SELECT id, name, secret FROM %s WHERE id=$1", appsTable)

	var app models.App
	if err := a.db.Get(&app, query, appID); err != nil {
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
