package http

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"avito-tenders/internal/api/middlewares"
)

const (
	tenderIDPathParam = "tenderId"
	bidIDPathParam    = "bidId"
	versionPathParam  = "version"
)

func (h *Handlers) MapBidsRoutes(r chi.Router, mw *middlewares.Manager) {
	r.Route("/bids", func(r chi.Router) {
		r.Post("/new", h.CreateBid)
		r.Get("/my", middlewares.Conveyor(h.GetMyBids, mw.UserExistsMiddleware, mw.PaginationMiddleware))
		r.Get(fmt.Sprintf("/{%s}/status", bidIDPathParam), middlewares.Conveyor(h.GetBidStatus, mw.UserExistsMiddleware))
		r.Put(fmt.Sprintf("/{%s}/status", bidIDPathParam), middlewares.Conveyor(h.UpdateBidStatus, mw.UserExistsMiddleware))
		r.Get(fmt.Sprintf("/{%s}/list", tenderIDPathParam), middlewares.Conveyor(h.FindBidsByTender, mw.UserExistsMiddleware, mw.PaginationMiddleware))
		r.Patch(fmt.Sprintf("/{%s}/edit", bidIDPathParam), middlewares.Conveyor(h.EditBid, mw.UserExistsMiddleware))
		r.Put(fmt.Sprintf("/{%s}/rollback/{%s}", bidIDPathParam, versionPathParam), middlewares.Conveyor(h.Rollback, mw.UserExistsMiddleware))

		r.Put(fmt.Sprintf("/{%s}/submit_decision", bidIDPathParam), middlewares.Conveyor(h.SubmitDecision, mw.UserExistsMiddleware))

		r.Put(fmt.Sprintf("/{%s}/feedback", bidIDPathParam), middlewares.Conveyor(h.SendFeedback, mw.UserExistsMiddleware))
		r.Get(fmt.Sprintf("/{%s}/reviews", tenderIDPathParam), middlewares.Conveyor(h.FindReviewsByTender, mw.PaginationMiddleware))
	})
}
