package middlewares

import (
	"context"
	"net/http"

	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/fwcontext"
	"avito-tenders/pkg/queryparams"
)

func (mw *Manager) PaginationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := queryparams.ParseQueryPagination(r.URL.Query())
		if err != nil {
			apperror.SendError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), fwcontext.PaginationCtxKey, pagination)

		next(w, r.WithContext(ctx))
	}
}
