package store

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	"github.com/ToraNoDora/little-sso/sso/internal/src/store/cache"
	p "github.com/ToraNoDora/little-sso/sso/internal/src/store/postgres"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Store struct {
	Cache      cache.Cache
	Repository *Repository
}

type Repository struct {
	Auth
	App
	User
	Permission
	Group
	Role
}

func NewStore(db *sqlx.DB, ch *cache.CacheStorage) *Store {
	return &Store{
		Cache: ch.Cache,
		Repository: &Repository{
			Auth:       p.NewAuthPostgres(db),
			App:        p.NewAppPostgres(db),
			User:       p.NewUserPostgres(db),
			Permission: p.NewPermissionPostgres(db),
			Group:      p.NewGroupPostgres(db),
			Role:       p.NewRolePostgres(db),
		},
	}
}

type Auth interface {
	CreateUser(ctx context.Context, username string, email string, passHash []byte) (string, error)
}

type App interface {
	GetApp(ctx context.Context, appID string) (models.App, error)
}

type Permission interface {
	GetUserPermissions(ctx context.Context, userID string) ([]models.Permission, error)
	AddPermission(ctx context.Context, email string, groupID int) (string, error)
	AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error)
	RemovePermission(ctx context.Context, email string, groupID int) (bool, error)
}

type Group interface {
	Create(ctx context.Context, name string, appID string) (int, error)
	Delete(ctx context.Context, groupID int) error
}

type Role interface {
	Create(ctx context.Context, name string, desc string) (int, error)
	AppointToGroup(ctx context.Context, groupID, roleID int) (int, error)
	Delete(ctx context.Context, roleID int) error
}

type User interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID string, appID string) (bool, error)
}
