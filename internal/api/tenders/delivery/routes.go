package delivery

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) MapTendersRoutes(r chi.Router) {
	r.Route("/tenders", func(r chi.Router) {
		r.Post("/new", h.CreateTender)
		r.Get("/my", h.GetMyTenders)
	})
}
