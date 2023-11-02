package server

import (
	"context"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/shimmy8/image-previewer/internal/config"
)

type Server struct {
	config *config.HttpConfig
	app    *app.App
}

func New(config *config.HttpConfig, app *app.App) *Server {
	return &Server{config: config, app: app}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
