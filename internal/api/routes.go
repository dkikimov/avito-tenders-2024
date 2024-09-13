package api

import (
	"log/slog"
	"net/http"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/v2/context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	bidsHttp "avito-tenders/internal/api/bids/delivery/http"
	bidsRepo "avito-tenders/internal/api/bids/repository"
	bidsUsecase "avito-tenders/internal/api/bids/usecase"
	empRepo "avito-tenders/internal/api/employee/repository"
	"avito-tenders/internal/api/middlewares"
	orgRepo "avito-tenders/internal/api/organization/repository"
	tendersHttp "avito-tenders/internal/api/tenders/delivery/http"
	tendersRepo "avito-tenders/internal/api/tenders/repository"
	tendersUsecase "avito-tenders/internal/api/tenders/usecase"
	"avito-tenders/pkg/backend"
)

const groupAPI = "/api"

func InitAPIRoutes(b backend.Backend) (chi.Router, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	tendersRepository := tendersRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)
	organizationRepository := orgRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)
	bidsRepository := bidsRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)
	empRepository := empRepo.NewRepository(b.DB, trmsqlx.DefaultCtxGetter)

	trManager := manager.Must(trmsqlx.NewDefaultFactory(b.DB), manager.WithCtxManager(trmcontext.DefaultManager))

	tendersUC := tendersUsecase.NewUsecase(tendersUsecase.Opts{
		Repo:      tendersRepository,
		OrgRepo:   organizationRepository,
		TrManager: trManager,
		EmpRepo:   empRepository,
	})
	bidsUC := bidsUsecase.NewUsecase(bidsUsecase.Opts{
		Repo:       bidsRepository,
		OrgRepo:    organizationRepository,
		EmpRepo:    empRepository,
		TenderRepo: tendersRepository,
		TrManager:  trManager,
	})

	mwManager := middlewares.NewManager(empRepository)

	tenderHandlers := tendersHttp.NewHandlers(tendersUC)
	bidsHandlers := bidsHttp.NewHandlers(bidsUC)

	r.Route(groupAPI, func(r chi.Router) {
		tenderHandlers.MapTendersRoutes(r, mwManager)
		bidsHandlers.MapBidsRoutes(r, mwManager)
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			err := b.DB.PingContext(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
	})

	slog.Info("API routes initialized")

	return r, nil
}
