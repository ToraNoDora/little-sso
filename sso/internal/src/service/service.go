package service

import (
	"context"
	"time"

	"golang.org/x/exp/slog"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	srv "github.com/ToraNoDora/little-sso/sso/internal/src/service/services"
	s "github.com/ToraNoDora/little-sso/sso/internal/src/store"
)

type Service struct {
	Authorization
	Permission
	User
}

func NewService(s *s.Store, log *slog.Logger, tokenTTL time.Duration) *Service {
	st := s.Repository
	return &Service{
		Authorization: srv.NewAuthService(log, st.Auth, st.User, st.App, tokenTTL, s.Cache),
		User:          srv.NewUserService(log, st.User),
		Permission:    srv.NewPermissionService(log, st.Permission),
	}
}

// Interfaces of service
type Authorization interface {
	SignIn(ctx context.Context, email string, password string, appID string) (string, error)
	SignUp(ctx context.Context, username string, email string, password string) (string, error)
}

type Permission interface {
	AddPermission(ctx context.Context, email string, groupID int) (string, error)
	AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error)
	RemovePermission(ctx context.Context, email string, groupID int) (bool, error)
}

type User interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID string, appID string) (bool, error)
}
