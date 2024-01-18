package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
)

type PermissionPostgres struct {
	db *sqlx.DB
}

func NewPermissionPostgres(db *sqlx.DB) *PermissionPostgres {
	return &PermissionPostgres{db: db}
}

func (p *PermissionPostgres) GetUserPermissions(ctx context.Context, userID string) ([]models.Permission, error) {
	const op = "storage.postgres.GetUserPermissions"

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", permissionsTable)

	var permissions []models.Permission
	err := p.db.SelectContext(ctx, &permissions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return permissions, nil
}

func (p *PermissionPostgres) AddPermission(ctx context.Context, email string, groupID int) (string, error) {
	const op = "storage.postgres.AddPermission"

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)

	var user models.User
	if err := p.db.Get(&user, query, email); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var prms []models.Permission
	query = fmt.Sprintf(
		"SELECT id, user_id, group_id FROM %s WHERE user_id = $1",
		permissionsTable,
	)

	if err := p.db.SelectContext(ctx, &prms, query, user.ID); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if checkPermissions(groupID, prms) {
		return "", fmt.Errorf("permission already exists")
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, group_id, add_flag) VALUES ($1, $2, $3) RETURNING id", permissionsTable)
	row := p.db.QueryRowContext(ctx, query, user.ID, groupID, true)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func checkPermissions(groupID int, prms []models.Permission) bool {
	for _, prm := range prms {
		if prm.GroupID == groupID {
			return true
		}
	}

	return false
}

func (p *PermissionPostgres) RemovePermission(ctx context.Context, email string, groupID int) (bool, error) {
	const op = "storage.postgres.RemovePermission"

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)

	var st bool
	var user models.User
	if err := p.db.Get(&user, query, email); err != nil {
		return st, fmt.Errorf("%s: %w", op, err)
	}

	query = fmt.Sprintf(
		`DELETE FROM %s WHERE user_id = $1 AND group_id = $2`,
		permissionsTable,
	)

	_, err := p.db.ExecContext(ctx, query, user.ID, groupID)
	if err != nil {
		return st, fmt.Errorf("%s: %w", op, err)
	} else {
		st = true
	}

	return st, nil
}

func (p *PermissionPostgres) AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error) {
	const op = "storage.postgres.AppointAsAdministrator"

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)

	var user models.User
	if err := p.db.Get(&user, query, email); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// checked admin role
	var admins []models.Admin
	query = fmt.Sprintf(
		"SELECT id FROM %s WHERE user_id = $1 AND app_id = $2",
		adminsTable,
	)

	if err := p.db.SelectContext(ctx, &admins, query, user.ID, appID); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(admins) > 0 {
		return "", fmt.Errorf("admin already exists")
	}

	// new admin
	query = fmt.Sprintf("INSERT INTO %s (user_id, app_id, is_admin) VALUES ($1, $2, $3) RETURNING id", adminsTable)
	row := p.db.QueryRowContext(ctx, query, user.ID, appID, is_admin)

	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
