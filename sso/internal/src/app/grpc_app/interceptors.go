package grpc_app

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/voidspooks/go-grpc-middleware/v2/interceptors"
	"github.com/voidspooks/go-grpc-middleware/v2/interceptors/auth"
	"github.com/voidspooks/go-grpc-middleware/v2/interceptors/logging"
)

var tokenInfoKey struct{}

// used by a middleware to authenticate requests
func authInterceptor(ctx context.Context) (context.Context, error) {
	const op = "grpcapp.authInterceptor"

	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	ctx = logging.InjectFields(ctx, logging.Fields{"auth.sub", userClaimFromToken(tokenInfo)})

	return context.WithValue(ctx, tokenInfoKey, tokenInfo), nil
}

func parseToken(token string) (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken(struct{}) string {
	return "foobar"
}

func signAnySkip(_ context.Context, c interceptors.CallMeta) bool {
	switch c.FullMethod() {
	case "/auth.Auth/SignUp", "/auth.Auth/SignIn":
		return false
	}

	return true
}
