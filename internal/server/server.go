package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/shimmy8/image-previewer/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	config *config.HttpConfig
	app    *app.App

	server *http.Server
	logger *zap.Logger
}

func New(config *config.HttpConfig, app *app.App, logger *zap.Logger) *Server {
	return &Server{config: config, app: app, server: nil, logger: logger}
}

func (s *Server) Start(ctx context.Context) error {
	handler := NewHandler(s.app, ctx, s.logger)

	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.config.Port),
		Handler: handler,
	}

	s.logger.Info("Starting server", zap.Int("port", s.config.Port))

	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server start error", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
