package dtos

import (
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/types"
)

type BidResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      entity.BidStatus  `json:"status"`
	AuthorType  entity.AuthorType `json:"authorType"`
	AuthorID    string            `json:"authorId"`
	Version     int               `json:"version"`
	CreatedAt   types.RFC3339Time `json:"createdAt"`
}

func NewBidResponse(bid entity.Bid) BidResponse {
	return BidResponse{
		ID:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.Status,
		AuthorType:  bid.AuthorType,
		AuthorID:    bid.AuthorID,
		Version:     bid.Version,
		CreatedAt:   types.RFCFromTime(bid.CreatedAt),
	}
}

func NewBidResponseList(bids []entity.Bid) []BidResponse {
	bidsResponseList := make([]BidResponse, 0, len(bids))

	for i := range bids {
		bidsResponseList = append(bidsResponseList, NewBidResponse(bids[i]))
	}

	return bidsResponseList
}
