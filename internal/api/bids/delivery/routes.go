package delivery

import (
	"github.com/go-chi/chi"
)

func MapBidsRoutes(r *chi.Mux, h HTTPHandlers) {
	r.Get("/ping", h.Ping)
}
