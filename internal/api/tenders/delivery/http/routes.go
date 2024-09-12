package http

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"avito-tenders/internal/api/middlewares"
)

const (
	tenderIdPathParam = "tenderId"
	versionPathParam  = "version"
)

func (h *Handlers) MapTendersRoutes(r chi.Router, mw *middlewares.Manager) {
	r.Route("/tenders", func(r chi.Router) {
		r.Get("/", middlewares.Conveyor(h.GetTenders, mw.PaginationMiddleware))
		r.Post("/new", h.CreateTender)
		r.Get("/my", middlewares.Conveyor(h.GetMyTenders, mw.UserExistsMiddleware))
		r.Get(fmt.Sprintf("/{%s}/status", tenderIdPathParam), h.GetTenderStatus)
		r.Put(fmt.Sprintf("/{%s}/status", tenderIdPathParam), middlewares.Conveyor(h.UpdateTenderStatus, mw.UserExistsMiddleware))
		r.Patch(fmt.Sprintf("/{%s}/edit", tenderIdPathParam), middlewares.Conveyor(h.UpdateTender, mw.UserExistsMiddleware))
		r.Put(fmt.Sprintf("/{%s}/rollback/{%s}", tenderIdPathParam, versionPathParam), middlewares.Conveyor(h.RollbackTender, mw.UserExistsMiddleware))
	})
}
