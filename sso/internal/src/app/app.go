package app

import (
	"time"

	"golang.org/x/exp/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // draver pq

	s "github.com/ToraNoDora/little-sso/sso/internal/src/store"
	"github.com/ToraNoDora/little-sso/sso/internal/src/store/cache"
	r "github.com/ToraNoDora/little-sso/sso/internal/src/store/cache/redis_cache"
	p "github.com/ToraNoDora/little-sso/sso/internal/src/store/postgres"

	"github.com/ToraNoDora/little-sso/sso/internal/src/service"

	ga "github.com/ToraNoDora/little-sso/sso/internal/src/app/grpc_app"
)

type App struct {
	db      *sqlx.DB
	GRPCSrv *ga.GRPCApp
}

func NewApp(
	log *slog.Logger,
	storeCfg p.Config,
	redisCfg r.Config,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	db, err := p.NewPostgresDB(storeCfg)
	if err != nil {
		panic(err)
	}

	rd := cache.NewCache(redisCfg)

	st := s.NewStore(db, rd)
	sr := service.NewService(st, log, tokenTTL)

	grpcApp := ga.NewGRPCApp(log, sr, grpcPort)

	return &App{
		db:      db,
		GRPCSrv: grpcApp,
	}
}

func (a *App) Stop(log *slog.Logger) error {
	if err := a.db.Close(); err != nil {
		log.Error(
			"failed to occured on db connection closed",
			err.Error(),
		)
		return err
	}

	if err := a.GRPCSrv.Stop(); err != nil {
		log.Error(
			"failed to stop gRPC server",
			err.Error(),
		)
		return err
	}

	return nil
}
