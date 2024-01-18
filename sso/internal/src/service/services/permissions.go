package services

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
)

type PermissionService struct {
	log         *slog.Logger
	prmProvider PermissionProvider
}

type PermissionProvider interface {
	AddPermission(ctx context.Context, email string, groupID int) (string, error)
	AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error)
	RemovePermission(ctx context.Context, email string, groupID int) (bool, error)
}

func NewPermissionService(
	log *slog.Logger,
	prmProvider PermissionProvider,
) *PermissionService {
	return &PermissionService{
		log:         log,
		prmProvider: prmProvider,
	}
}

func (p *PermissionService) AddPermission(ctx context.Context, email string, groupID int) (string, error) {
	const op = "service.permissions.AddPermission"

	log := p.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int("group_id", groupID),
	)

	id, err := p.prmProvider.AddPermission(ctx, email, groupID)
	if err != nil {
		log.Error("failed to add permission to user", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *PermissionService) AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error) {
	const op = "service.permissions.AppointAsAdministrator"

	log := p.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("app_id", appID),
	)

	id, err := p.prmProvider.AppointAsAdmin(ctx, email, appID, is_admin)
	if err != nil {
		log.Error("failed to appoint user as admin", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *PermissionService) RemovePermission(ctx context.Context, email string, groupID int) (bool, error) {
	const op = "service.permissions.RemovePermission"

	log := p.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int("groupID", groupID),
	)

	st, err := p.prmProvider.RemovePermission(ctx, email, groupID)
	if err != nil {
		log.Error("failed to remove permission at user", err.Error())
		return st, fmt.Errorf("%s: %w", op, err)
	}

	return st, nil
}
