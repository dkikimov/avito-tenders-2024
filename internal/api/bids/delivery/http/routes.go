package http

import (
	"github.com/go-chi/chi/v5"
)

func MapBidsRoutes(r *chi.Mux, h HTTPHandlers) {
	r.Get("/ping", h.Ping)
}
