package http

import (
	"encoding/json"
	"io"
	"net/http"

	"avito-tenders/internal/api/bids"
	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/pkg/apperror"
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
