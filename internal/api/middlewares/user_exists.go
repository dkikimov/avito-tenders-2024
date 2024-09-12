package middlewares

import (
	"context"
	"net/http"

	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/fwcontext"
)

// UserExistsMiddleware checks if the user with the given ID exists.
func (mw *Manager) UserExistsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the user ID from the query parameter
		username := r.URL.Query().Get("username")
		if username == "" {
			// Handle the case where the userID is not provided
			apperror.SendError(w, apperror.Unauthorized(apperror.ErrUserEmpty))
			return
		}

		// Check if the user exists in the repository
		_, err := mw.empRepo.FindByUsername(r.Context(), username)
		if err != nil {
			apperror.SendError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), fwcontext.UsernameCtxKey, username)

		// Pass the request to the next handler if the user exists
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
