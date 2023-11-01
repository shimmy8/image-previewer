package main

import (
	"context"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/shimmy8/image-previewer/internal/config"
	"github.com/shimmy8/image-previewer/internal/server"
)

func main() {
	config := config.New()
	app := app.New(&config.Cache)
	server := server.New(&config.Http, app)

	ctx := context.Background()

	server.Start(ctx)
	// TODO
}
