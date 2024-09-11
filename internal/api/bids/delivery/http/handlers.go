package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"avito-tenders/internal/api/bids"
	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
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
	values := r.URL.Query()

	username := values.Get("username")
	if username == "" {
		apperror.SendError(w, apperror.Unauthorized(errors.New("username is required")))
		return
	}

	pagination, err := queryparams.ParsePagination(values)
	if err != nil {
		apperror.SendError(w, err)
	}

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
