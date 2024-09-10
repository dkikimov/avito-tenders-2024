package delivery

import (
	"fmt"

	"github.com/go-chi/chi/v5"
)

const tenderIdPathParam = "tenderId"

func (h *Handlers) MapTendersRoutes(r chi.Router) {
	r.Route("/tenders", func(r chi.Router) {
		r.Post("/new", h.CreateTender)
		r.Get("/my", h.GetMyTenders)
		r.Get(fmt.Sprintf("/{%s}/status", tenderIdPathParam), h.GetTenderStatus)
		r.Put(fmt.Sprintf("/{%s}/status", tenderIdPathParam), h.UpdateTenderStatus)
		r.Patch(fmt.Sprintf("/{%s}/edit", tenderIdPathParam), h.UpdateTender)
	})
}
