package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (u *UserPostgres) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.GetUser"

	query := fmt.Sprintf(
		"SELECT id, username, email, pass_hash, created_at FROM %s WHERE email=$1",
		usersTable,
	)

	var user models.User
	if err := u.db.Get(&user, query, email); err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	var prms []models.Permission
	query = fmt.Sprintf(
		"SELECT id, user_id, group_id FROM %s WHERE user_id = $1",
		permissionsTable,
	)

	if err := u.db.SelectContext(ctx, &prms, query, user.ID); err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	user.Permissions = prms

	return user, nil
}

func (u *UserPostgres) IsAdmin(ctx context.Context, userID string, appID string) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := fmt.Sprintf("SELECT id, is_admin FROM %s WHERE user_id=$1 AND app_id=$2", adminsTable)

	var admin models.Admin
	if err := u.db.Get(&admin, query, userID, appID); err != nil {
		return admin.IsAdmin, fmt.Errorf("%s: %w", op, err)
	}

	return admin.IsAdmin, nil
}
