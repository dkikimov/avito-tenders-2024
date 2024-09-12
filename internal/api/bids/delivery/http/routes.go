package http

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"avito-tenders/internal/api/middlewares"
)

const (
	tenderIdPathParam = "tenderId"
	bidIdPathParam    = "bidId"
)

func (h *Handlers) MapBidsRoutes(r chi.Router, mw *middlewares.Manager) {
	r.Route("/bids", func(r chi.Router) {
		r.Post("/new", h.CreateBid)
		r.Get("/my", middlewares.Conveyor(h.GetMyBids, mw.UserExistsMiddleware, mw.PaginationMiddleware))
		r.Get(fmt.Sprintf("/{%s}/status", bidIdPathParam), middlewares.Conveyor(h.GetBidStatus, mw.UserExistsMiddleware))
		r.Put(fmt.Sprintf("/{%s}/status", bidIdPathParam), middlewares.Conveyor(h.UpdateBidStatus, mw.UserExistsMiddleware))
		r.Get(fmt.Sprintf("/{%s}/list", tenderIdPathParam), middlewares.Conveyor(h.FindBidsByTender, mw.UserExistsMiddleware, mw.PaginationMiddleware))
		r.Patch(fmt.Sprintf("/{%s}/edit", bidIdPathParam), middlewares.Conveyor(h.EditBid, mw.UserExistsMiddleware))

		r.Put(fmt.Sprintf("/{%s}/submit_decision", bidIdPathParam), middlewares.Conveyor(h.SubmitDecision, mw.UserExistsMiddleware))
	})
}
