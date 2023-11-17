package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/shimmy8/image-previewer/internal/app"
	"go.uber.org/zap"
)

type Server struct {
	port int

	app *app.App

	server *http.Server
	logger *zap.Logger
}

func New(port int, app *app.App, logger *zap.Logger) *Server {
	return &Server{port: port, app: app, server: nil, logger: logger}
}

func (s *Server) Start(ctx context.Context) error {
	handler := NewHandler(ctx, s.app, s.logger)

	s.server = &http.Server{
		Addr:              ":" + strconv.Itoa(s.port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 1,
	}

	s.logger.Info("Starting server", zap.Int("port", s.port))

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
