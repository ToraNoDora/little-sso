package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/ToraNoDora/little-sso/sso/internal/src/app"
	"github.com/ToraNoDora/little-sso/sso/internal/src/config"
	slp "github.com/ToraNoDora/little-sso/sso/pkg/logger/handlers/slog_pretty"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/config.locale.yaml", "path to config")
}

const (
	envLocal = "locale"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := setupLogger(cfg.Env)
	log.Info(
		"starting app",
	)

	app := app.NewApp(log, cfg.Store, cfg.Redis, cfg.GRPC.Port, cfg.TokenTtl)

	go app.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)

	sign := <-stop

	log.Info(
		"stopping app",
		slog.String("signal", sign.String()),
	)

	if err := app.Stop(log); err != nil {
		log.Error(
			"failed to stop server",
			err.Error(),
		)
	}

	log.Info("app stopped!")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slp.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
