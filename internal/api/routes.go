package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	orgRepo "avito-tenders/internal/api/organization/repository"
	"avito-tenders/internal/api/tenders/delivery/http"
	tendersRepo "avito-tenders/internal/api/tenders/repository"
	"avito-tenders/internal/api/tenders/usecase"
	"avito-tenders/pkg/backend"
)

const groupAPI = "/api"

// InitAPIRoutes initializes router for the webserver.
// Swagger spec:
// @schemes     https
// @host        localhost:8080
// @BasePath    /api
// @title       Avito Tenders API
// @version     1.0.
func InitAPIRoutes(b backend.Backend) (chi.Router, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	tendersRepository := tendersRepo.NewRepository(b.DB)
	organizationRepository := orgRepo.NewRepository(b.DB)

	tendersUC := usecase.NewUseCase(tendersRepository, organizationRepository)

	tenderHandlers := http.NewHandlers(tendersUC)

	r.Route(groupAPI, func(r chi.Router) {
		tenderHandlers.MapTendersRoutes(r)
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	slog.Info("API routes initialized")

	return r, nil
}
