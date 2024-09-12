package middlewares

import "net/http"

// Conveyor is helper function that applies all middlewares to handlers.
// Middlewares apply in the direct order.
func Conveyor(handler http.HandlerFunc, middlewaresList ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next := handler

		for i := len(middlewaresList) - 1; i >= 0; i-- {
			next = middlewaresList[i](next)
		}

		next.ServeHTTP(w, r)
	}
}
