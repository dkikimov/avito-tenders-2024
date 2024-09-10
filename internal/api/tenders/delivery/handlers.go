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
	"avito-tenders/internal/entity"
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

	if err := tender.Validate(); err != nil {
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

func (h *Handlers) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, tenderIdPathParam)
	if tenderId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	urlQuery := r.URL.Query()

	req := entities.EditTenderStatusRequest{
		Status:   entity.TenderStatus(urlQuery.Get("status")),
		Username: urlQuery.Get("username"),
	}

	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.EditStatus(r.Context(), tenderId, req)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) UpdateTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, tenderIdPathParam)
	if tenderId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	urlQuery := r.URL.Query()
	username := urlQuery.Get("username")
	if len(username) == 0 {
		apperror.SendError(w, apperror.Unauthorized(apperror.ErrUserDoesNotExist))
		return
	}

	// Parse body to EditTender.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	var edit entities.EditTender
	if err := json.Unmarshal(body, &edit); err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	if valid, err := govalidator.ValidateStruct(&edit); !valid || err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Edit(r.Context(), tenderId, entities.EditTenderRequest{
		EditTender: edit,
		Username:   username,
	})
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) RollbackTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, tenderIdPathParam)
	if tenderId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	version := chi.URLParam(r, versionPathParam)
	if version == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("version is not specified")))
		return
	}

	urlQuery := r.URL.Query()
	username := urlQuery.Get("username")
	if len(username) == 0 {
		apperror.SendError(w, apperror.Unauthorized(apperror.ErrUserDoesNotExist))
		return
	}

	request := entities.RollbackTenderRequest{
		Username: username,
		Version:  version,
	}
	if err := request.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Rollback(r.Context(), tenderId, request)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}
