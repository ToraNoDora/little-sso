package permission

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/ToraNoDora/little-sso-protos/gen/go/sso"

	srv "github.com/ToraNoDora/little-sso/sso/internal/src/service/services"
)

const (
	emptyValueInt = 0
	emptyValueStr = ""
)

type Permission interface {
	AddPermission(ctx context.Context, userID string, groupID int) (string, error)
	AppointAsAdmin(ctx context.Context, email string, appID string, is_admin bool) (string, error)
	RemovePermission(ctx context.Context, email string, groupID int) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedPermissionServer
	permission Permission
}

func NewPermissionServer(gRPC *grpc.Server, permission Permission) {
	ssov1.RegisterPermissionServer(gRPC, &serverApi{permission: permission})
}

func (s *serverApi) AddPermission(ctx context.Context, req *ssov1.AddPermissionRequest) (*ssov1.AddPermissionResponse, error) {
	if err := validatePermissionReq(req); err != nil {
		return nil, err
	}

	id, err := s.permission.AddPermission(ctx, req.GetEmail(), int(req.GetGroupId()))
	if err != nil {
		if errors.Is(err, srv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to add permission")
	}

	return &ssov1.AddPermissionResponse{
		PermissionId: id,
	}, nil
}

func validatePermissionReq(req *ssov1.AddPermissionRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetGroupId() == emptyValueInt {
		return status.Error(codes.InvalidArgument, "group_id is required")
	}

	return nil
}

func (s *serverApi) AppointAsAdministrator(ctx context.Context, req *ssov1.AppointAsAdministratorRequest) (*ssov1.AppointAsAdministratorResponse, error) {
	if err := validateAppointAsAdminReq(req); err != nil {
		return nil, err
	}

	adminID, err := s.permission.AppointAsAdmin(ctx, req.GetEmail(), req.GetAppId(), req.GetIsAdmin())
	if err != nil {
		if errors.Is(err, srv.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.AppointAsAdministratorResponse{
		AdminId: adminID,
	}, nil
}

func validateAppointAsAdminReq(req *ssov1.AppointAsAdministratorRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetAppId() == emptyValueStr {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func (s *serverApi) RemovePermission(ctx context.Context, req *ssov1.RemovePermissionRequest) (*ssov1.RemovePermissionResponse, error) {
	if err := validateRemovePermissionReq(req); err != nil {
		return nil, err
	}

	st, err := s.permission.RemovePermission(ctx, req.GetEmail(), int(req.GetGroupId()))
	if err != nil {
		if errors.Is(err, srv.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RemovePermissionResponse{
		Status: st,
	}, nil
}

func validateRemovePermissionReq(req *ssov1.RemovePermissionRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetGroupId() == emptyValueInt {
		return status.Error(codes.InvalidArgument, "group_id is required")
	}

	return nil
}
