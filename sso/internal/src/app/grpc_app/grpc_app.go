package grpc_app

import (
	"fmt"
	"net"

	"golang.org/x/exp/slog"

	"github.com/voidspooks/go-grpc-middleware/v2/interceptors/auth"
	"github.com/voidspooks/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"

	g "github.com/ToraNoDora/little-sso/sso/internal/src/grpc"
	"github.com/ToraNoDora/little-sso/sso/internal/src/service"
)

type GRPCApp struct {
	log        *slog.Logger
	port       int
	gRPCServer *grpc.Server
}

func NewGRPCApp(log *slog.Logger, s *service.Service, port int) *GRPCApp {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(
				auth.UnaryServerInterceptor(authInterceptor),
				selector.MatchFunc(signAnySkip),
			),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(
				auth.StreamServerInterceptor(authInterceptor),
				selector.MatchFunc(signAnySkip),
			),
		),
	)

	g.NewGRPCService(gRPCServer, s.Authorization, s.Permission, s.User)

	return &GRPCApp{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (ga *GRPCApp) MustRun() {
	if err := ga.Run(); err != nil {
		panic(err)
	}
}

func (ga *GRPCApp) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", ga.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ga.log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := ga.gRPCServer.Serve(l); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (ga *GRPCApp) Stop() error {
	const op = "grpcapp.Stop"

	ga.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", ga.port))

	ga.gRPCServer.GracefulStop()

	return nil
}
