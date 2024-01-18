package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ssov1 "github.com/ToraNoDora/little-sso-protos/gen/go/sso"

	"github.com/ToraNoDora/little-sso/sso/internal/src/config"
)

const (
	grpcHost = "localhost"
)

var Cfg = config.MustLoadByPath(configPath())

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../configs/config.locale.yaml"
}

type Suite struct {
	*testing.T
	AuthClient ssov1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	ctx, cancelCtx := context.WithTimeout(context.Background(), Cfg.GRPC.Timeout)

	t.Cleanup(
		func() {
			t.Helper()
			cancelCtx()
		},
	)

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(Cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(gRPCPort int) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(gRPCPort))
}
