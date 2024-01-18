package grpc

import (
	"google.golang.org/grpc"

	authgrpc "github.com/ToraNoDora/little-sso/sso/internal/src/grpc/auth"
	prmgrpc "github.com/ToraNoDora/little-sso/sso/internal/src/grpc/permission"
	usergrpc "github.com/ToraNoDora/little-sso/sso/internal/src/grpc/user"
)

func NewGRPCService(
	gRPC *grpc.Server,
	authSrv authgrpc.Auth,
	prmSrv prmgrpc.Permission,
	userSrv usergrpc.User,
) {
	authgrpc.NewAuthServer(gRPC, authSrv)
	prmgrpc.NewPermissionServer(gRPC, prmSrv)
	usergrpc.NewUserServer(gRPC, userSrv)
}
