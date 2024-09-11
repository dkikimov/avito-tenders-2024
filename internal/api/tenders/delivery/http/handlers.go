package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/invopop/validation"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
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

	var tender dtos.CreateTenderRequest
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

	pagination, err := queryparams.ParsePagination(values)
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

func (h *Handlers) GetTenders(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	// Parse pagination.
	pagination, err := queryparams.ParsePagination(values)
	if err != nil {
		apperror.SendError(w, err)
		slog.Error("couldn't parse pagination", "error", err)
		return
	}

	// Parse service types.
	serviceTypesStrings := values["service_type"]
	serviceTypeList := make([]entity.ServiceType, 0, len(serviceTypesStrings))
	for _, serviceTypeString := range serviceTypesStrings {
		serviceType := entity.ServiceType(serviceTypeString)

		if err := validation.Validate(serviceType, serviceType.ValidationRule()); err != nil {
			apperror.SendError(w, apperror.BadRequest(err))
			return
		}

		serviceTypeList = append(serviceTypeList, serviceType)
	}

	// Getting all tenders with filter.
	createdTender, err := h.uc.GetAll(r.Context(), tenders.TenderFilter{
		ServiceTypes: serviceTypeList,
	}, pagination)
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

	request := dtos.TenderStatus{
		Username: r.URL.Query().Get("username"),
	}

	tender, err := h.uc.GetTenderStatus(r.Context(), tenderId, request)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dtos.TenderStatusResponse{Status: tender.Status}); err != nil {
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

	req := dtos.EditTenderStatusRequest{
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

	var edit dtos.EditTender
	if err := json.Unmarshal(body, &edit); err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	if valid, err := govalidator.ValidateStruct(&edit); !valid || err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Edit(r.Context(), tenderId, dtos.EditTenderRequest{
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

	request := dtos.RollbackTenderRequest{
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
