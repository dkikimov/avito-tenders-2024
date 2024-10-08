package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/invopop/validation"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/fwcontext"
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
	username := fwcontext.GetUsername(r.Context())
	pagination := fwcontext.GetPagination(r.Context())

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

	pagination := fwcontext.GetPagination(r.Context())

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
	tenderID := chi.URLParam(r, tenderIDPathParam)
	if tenderID == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	request := dtos.TenderStatus{
		Username: r.URL.Query().Get("username"),
	}

	tender, err := h.uc.GetTenderStatus(r.Context(), tenderID, request)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write([]byte(tender.Status)); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderID := chi.URLParam(r, tenderIDPathParam)
	if tenderID == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	urlQuery := r.URL.Query()

	req := dtos.EditTenderStatusRequest{
		Status:   entity.TenderStatus(urlQuery.Get("status")),
		Username: fwcontext.GetUsername(r.Context()),
	}

	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.EditStatus(r.Context(), tenderID, req)
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
	tenderID := chi.URLParam(r, tenderIDPathParam)
	if tenderID == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	username := fwcontext.GetUsername(r.Context())

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

	req := dtos.EditTenderRequest{
		EditTender: edit,
		Username:   username,
	}
	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Edit(r.Context(), tenderID, req)
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
	tenderID := chi.URLParam(r, tenderIDPathParam)
	if tenderID == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	version := chi.URLParam(r, versionPathParam)
	if version == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("version is not specified")))
		return
	}

	username := fwcontext.GetUsername(r.Context())

	intVersion, err := strconv.Atoi(version)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(errors.New("version is not a number")))
		return
	}

	request := dtos.RollbackTenderRequest{
		Username: username,
		Version:  intVersion,
	}
	if err := request.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Rollback(r.Context(), tenderID, request)
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
