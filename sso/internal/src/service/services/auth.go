package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"

	"github.com/ToraNoDora/little-sso/sso/pkg/logger/sl"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	"github.com/ToraNoDora/little-sso/sso/internal/src/lib/hash"
	"github.com/ToraNoDora/little-sso/sso/internal/src/lib/jwt"
	"github.com/ToraNoDora/little-sso/sso/internal/src/store"
	"github.com/ToraNoDora/little-sso/sso/internal/src/store/cache"
)

type AuthService struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
	cache       cache.Cache
}

type UserSaver interface {
	CreateUser(ctx context.Context, username string, email string, password []byte) (string, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID string) (models.App, error)
}

func NewAuthService(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
	cache cache.Cache,
) *AuthService {
	return &AuthService{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
		log:         log,
		cache:       cache,
	}
}

func (a *AuthService) SignUp(ctx context.Context, username string, email string, password string) (string, error) {
	const op = "service.auth.SignUp"

	log := a.log.With(slog.String("op", op))
	log.Info(
		"registering new user",
		slog.String("email", email),
	)

	passHash, err := hash.HashPassword(password)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.CreateUser(ctx, username, email, passHash)
	if err != nil {
		if errors.Is(err, store.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Login checks if user with given credentials exists in the system.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *AuthService) SignIn(ctx context.Context, email string, password string, appID string) (string, error) {
	const op = "auth.service.SignIn"

	log := a.log.With(slog.String("op", op))
	log.Info(
		"attempting to login user",
		slog.String("email", email),
	)

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := hash.VerifyPassword(password, user.PassHash); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	userPrms := user.Permissions
	hashString, err := hash.HashingPermissions(userPrms)
	if err != nil {
		log.Error("failed to hasher token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// caching permissions
	_, err = a.cachedPermissions(userPrms, hashString)
	if err != nil {
		log.Error("failed to caching user permissions", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL, hashString)
	if err != nil {
		log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) cachedPermissions(prms []models.Permission, hashString string) (bool, error) {
	permissionsJSON, err := json.Marshal(prms)
	result, err := a.cache.Set(hashString, permissionsJSON)
	if err != nil {
		return false, err
	}

	return result, nil
}
