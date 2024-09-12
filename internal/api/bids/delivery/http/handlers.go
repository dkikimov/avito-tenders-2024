package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"avito-tenders/internal/api/bids"
	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/fwcontext"
)

type Handlers struct {
	uc bids.Usecase
}

func NewHandlers(uc bids.Usecase) *Handlers {
	return &Handlers{uc: uc}
}

func (h *Handlers) CreateBid(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	var bid dtos.CreateBidRequest
	if err := json.Unmarshal(body, &bid); err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	if err := bid.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	createdBid, err := h.uc.Create(r.Context(), bid)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(createdBid); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) GetMyBids(w http.ResponseWriter, r *http.Request) {
	username := fwcontext.GetUsername(r.Context())
	pagination := fwcontext.GetPagination(r.Context())

	bidsList, err := h.uc.FindByUsername(r.Context(), username, pagination)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bidsList); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) FindBidsByTender(w http.ResponseWriter, r *http.Request) {
	tenderId := chi.URLParam(r, tenderIdPathParam)
	if tenderId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("tender id is not specified")))
		return
	}

	username := fwcontext.GetUsername(r.Context())
	pagination := fwcontext.GetPagination(r.Context())

	req := dtos.FindByTenderIdRequest{
		TenderId:   tenderId,
		Username:   username,
		Pagination: pagination,
	}
	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	bidsList, err := h.uc.FindByTenderId(r.Context(), req)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bidsList); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	username := fwcontext.GetUsername(r.Context())

	bidId := chi.URLParam(r, bidIdPathParam)
	if bidId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("bidId is not specified")))
		return
	}

	status, err := h.uc.GetStatusById(r.Context(), bidId, username)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dtos.BidStatusResponse{BidStatus: status}); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	username := fwcontext.GetUsername(r.Context())

	bidId := chi.URLParam(r, bidIdPathParam)
	if bidId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("bidId is not specified")))
		return
	}

	statusString := r.URL.Query().Get("status")
	if statusString == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("status is not specified")))
		return
	}

	req := dtos.UpdateStatusRequest{
		BidId:    bidId,
		Status:   entity.BidStatus(statusString),
		Username: username,
	}

	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	updatedBid, err := h.uc.UpdateStatusById(r.Context(), req)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedBid); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) EditBid(w http.ResponseWriter, r *http.Request) {
	username := fwcontext.GetUsername(r.Context())

	bidId := chi.URLParam(r, bidIdPathParam)
	if bidId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("bidId is not specified")))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	var bidBody dtos.EditBidBody
	if err := json.Unmarshal(body, &bidBody); err != nil {
		apperror.SendError(w, apperror.BadRequest(apperror.ErrInvalidInput))
		return
	}

	req := dtos.EditBidRequest{
		BidId:       bidId,
		Username:    username,
		EditBidBody: bidBody,
	}

	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	updatedBid, err := h.uc.Edit(r.Context(), req)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedBid); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) SubmitDecision(w http.ResponseWriter, r *http.Request) {
	username := fwcontext.GetUsername(r.Context())

	bidId := chi.URLParam(r, bidIdPathParam)
	if bidId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("bidId is not specified")))
		return
	}

	decision := r.URL.Query().Get("decision")

	req := dtos.SubmitDecisionRequest{
		BidId:    bidId,
		Decision: entity.BidDecision(decision),
		Username: username,
	}

	if err := req.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	updatedBid, err := h.uc.SubmitDecision(r.Context(), req)
	if err != nil {
		apperror.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedBid); err != nil {
		apperror.SendError(w, apperror.InternalServerError(err))
	}
}

func (h *Handlers) Rollback(w http.ResponseWriter, r *http.Request) {
	bidId := chi.URLParam(r, bidIdPathParam)
	if bidId == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("bidId is not specified")))
		return
	}

	version := chi.URLParam(r, versionPathParam)
	if version == "" {
		apperror.SendError(w, apperror.BadRequest(errors.New("version is not specified")))
		return
	}

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		apperror.SendError(w, apperror.BadRequest(errors.New("version is not a number")))
		return
	}

	username := fwcontext.GetUsername(r.Context())

	request := dtos.RollbackRequest{
		BidId:    bidId,
		Version:  versionInt,
		Username: username,
	}
	if err := request.Validate(); err != nil {
		apperror.SendError(w, apperror.BadRequest(err))
		return
	}

	tender, err := h.uc.Rollback(r.Context(), request)
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
