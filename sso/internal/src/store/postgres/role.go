package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RolePostgres struct {
	db *sqlx.DB
}

func NewRolePostgres(db *sqlx.DB) *RolePostgres {
	return &RolePostgres{db: db}
}

func (r *RolePostgres) Create(ctx context.Context, name string, desc string) (int, error) {
	const op = "storage.postgres.CreateRole"

	query := fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", rolesTable)
	row := r.db.QueryRowContext(ctx, query, name, desc)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *RolePostgres) AppointToGroup(ctx context.Context, groupID, roleID int) (int, error) {
	const op = "storage.postgres.AppointToGroup"

	// TODO check roles

	query := fmt.Sprintf("INSERT INTO %s (group_id, role_id) VALUES ($1, $2) RETURNING id", groupsRolesTable)
	row := r.db.QueryRowContext(ctx, query, groupID, roleID)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *RolePostgres) Delete(ctx context.Context, roleID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE role_id = $1`,
		groupsRolesTable,
	)
	_, err := r.db.ExecContext(ctx, query, roleID)

	query = fmt.Sprintf(
		`DELETE FROM %s WHERE id = $1`,
		rolesTable,
	)
	_, err = r.db.ExecContext(ctx, query, roleID)

	return err
}
