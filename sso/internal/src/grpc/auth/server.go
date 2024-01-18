package auth

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
	emptyValue = ""
)

type Auth interface {
	SignIn(ctx context.Context, email string, password string, appID string) (token string, err error)
	SignUp(ctx context.Context, username string, email string, password string) (userID string, err error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func NewAuthServer(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) SignIn(ctx context.Context, req *ssov1.SignInRequest) (*ssov1.SignInResponse, error) {
	if err := validateSignInReq(req); err != nil {
		return nil, err
	}

	token, err := s.auth.SignIn(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, srv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.SignInResponse{
		Token: token,
	}, nil
}

func validateSignInReq(req *ssov1.SignInRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func (s *serverApi) SignUp(ctx context.Context, req *ssov1.SignUpRequest) (*ssov1.SignUpResponse, error) {
	if err := validateSignUpReq(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.SignUp(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, srv.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.SignUpResponse{
		UserId: userID,
	}, nil
}

func validateSignUpReq(req *ssov1.SignUpRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
