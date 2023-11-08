package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/shimmy8/image-previewer/internal/config"
	"github.com/shimmy8/image-previewer/internal/server"
	"go.uber.org/zap"
)

func main() {
	// create logger
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	logger.Info("Starting app")

	config, err := config.New()
	if err != nil {
		logger.Named("main").Error("Config parse error", zap.Error(err))
	}

	app := app.New(config.Cache, logger.Named("app"))
	server := server.New(config.Http, app, logger.Named("server"))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		defer cancel()
		err := server.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", zap.Error(err))
		}
	}()

	if err := server.Start(ctx); err != nil {
		cancel()
	}

}
