package server

import (
	"net/http"

	"github.com/shimmy8/image-previewer/internal/app"
)

type Handler struct {
	app app.App
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) handleResizeRequest(w http.ResponseWriter, r *http.Request) error {
	//TODO
	return nil
}
