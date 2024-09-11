package http

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) MapBidsRoutes(r chi.Router) {
	r.Route("/bids", func(r chi.Router) {
		r.Post("/new", h.CreateBid)
	})
}
