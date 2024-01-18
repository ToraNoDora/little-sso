package services

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/ToraNoDora/little-sso/sso/pkg/logger/sl"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	"github.com/ToraNoDora/little-sso/sso/internal/src/store"
)

type UserService struct {
	log         *slog.Logger
	usrProvider UserProvider
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID string, appID string) (bool, error)
}

func NewUserService(
	log *slog.Logger,
	userProvider UserProvider,
) *UserService {
	return &UserService{
		usrProvider: userProvider,
		log:         log,
	}
}

func (a *UserService) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "service.user.GetUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		log.Error("failed to get user", err.Error())
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *UserService) IsAdmin(ctx context.Context, userID string, appID string) (bool, error) {
	const op = "service.user.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID, appID)
	if err != nil {
		if errors.Is(err, store.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check if user is admin", err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
