package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type GroupPostgres struct {
	db *sqlx.DB
}

func NewGroupPostgres(db *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{db: db}
}

func (g *GroupPostgres) Create(ctx context.Context, name string, appID string) (int, error) {
	const op = "storage.postgres.CreateGroup"

	query := fmt.Sprintf("INSERT INTO %s (name, app_id) VALUES ($1, $2) RETURNING id", groupsTable)
	row := g.db.QueryRowContext(ctx, query, name, appID)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (g *GroupPostgres) Delete(ctx context.Context, groupID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE group_id = $1`,
		groupsRolesTable,
	)

	_, err := g.db.ExecContext(ctx, query, groupID)

	query = fmt.Sprintf(
		`DELETE FROM %s WHERE group_id = $1`,
		permissionsTable,
	)

	_, err = g.db.ExecContext(ctx, query, groupID)

	query = fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1`,
		groupsTable,
	)

	_, err = g.db.ExecContext(ctx, query, groupID)

	return err
}
