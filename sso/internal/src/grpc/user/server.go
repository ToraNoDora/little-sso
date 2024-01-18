package user

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ssov1 "github.com/ToraNoDora/little-sso-protos/gen/go/sso"

	"github.com/ToraNoDora/little-sso/sso/internal/src/domain/models"
	srv "github.com/ToraNoDora/little-sso/sso/internal/src/service/services"
)

const (
	emptyValue = ""
)

type User interface {
	GetUser(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID string, appID string) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedUserServer
	user User
}

func NewUserServer(gRPC *grpc.Server, user User) {
	ssov1.RegisterUserServer(gRPC, &serverApi{user: user})
}

func (s *serverApi) GetUser(ctx context.Context, req *ssov1.GetUserRequest) (*ssov1.GetUserResponse, error) {
	if err := validateGetUserReq(req); err != nil {
		return nil, err
	}

	user, err := s.user.GetUser(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, srv.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.GetUserResponse{
		UserId:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}

func validateGetUserReq(req *ssov1.GetUserRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	return nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdminReq(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.user.IsAdmin(ctx, req.GetUserId(), req.GetAppId())
	if err != nil {
		if errors.Is(err, srv.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateIsAdminReq(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}
