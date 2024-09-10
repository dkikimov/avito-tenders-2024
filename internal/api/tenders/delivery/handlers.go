package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/query"
)

type Handlers struct {
	uc tenders.Usecase
}

func NewHandlers(uc tenders.Usecase) *Handlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) CreateTender(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	var tender entities.CreateTenderRequest
	if err := json.Unmarshal(body, &tender); err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	if _, err := govalidator.ValidateStruct(&tender); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	createdTender, err := h.uc.Create(r.Context(), tender)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(createdTender); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	username := values.Get("username")
	if len(username) == 0 {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	pagination, err := query.ParsePagination(values)
	if err != nil {
		apperror.SendError(w, err)
		slog.Error("couldn't parse pagination", "error", err)
		return
	}

	createdTender, err := h.uc.FindByUsername(r.Context(), username, pagination)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(createdTender); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, tenderIdPathParam)
	if tenderId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	tender, err := h.uc.FindById(r.Context(), tenderId)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entities.TenderStatusResponse{Status: tender.Status}); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}
