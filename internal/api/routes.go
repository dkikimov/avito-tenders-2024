package api

import (
	"log/slog"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	bidsHttp "avito-tenders/internal/api/bids/delivery/http"
	bidsRepo "avito-tenders/internal/api/bids/repository"
	bidsUsecase "avito-tenders/internal/api/bids/usecase"
	empRepo "avito-tenders/internal/api/employee/repository"
	orgRepo "avito-tenders/internal/api/organization/repository"
	tendersHttp "avito-tenders/internal/api/tenders/delivery/http"
	tendersRepo "avito-tenders/internal/api/tenders/repository"
	tendersUsecase "avito-tenders/internal/api/tenders/usecase"
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
	organizationRepository := orgRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)
	bidsRepository := bidsRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)
	empRepository := empRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)

	trManager := manager.Must(trmsqlx.NewDefaultFactory(b.DB), manager.WithCtxManager(trmcontext.DefaultManager))

	tendersUC := tendersUsecase.NewUseCase(tendersRepository, organizationRepository)
	bidsUC := bidsUsecase.NewUsecase(bidsUsecase.Opts{
		Repo:      bidsRepository,
		OrgRepo:   organizationRepository,
		EmpRepo:   empRepository,
		TrManager: trManager,
	})

	tenderHandlers := tendersHttp.NewHandlers(tendersUC)
	bidsHandlers := bidsHttp.NewHandlers(bidsUC)

	r.Route(groupAPI, func(r chi.Router) {
		tenderHandlers.MapTendersRoutes(r)
		bidsHandlers.MapBidsRoutes(r)
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	slog.Info("API routes initialized")

	return r, nil
}
