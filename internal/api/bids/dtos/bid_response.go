package dtos

import (
	"time"

	"avito-tenders/internal/entity"
)

type BidResponse struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Status     entity.BidStatus  `json:"status"`
	AuthorType entity.AuthorType `json:"author_type"`
	AuthorId   string            `json:"author_id"`
	Version    int               `json:"version"`
	CreatedAt  time.Time         `json:"created_at"`
}

func NewBidResponse(bid entity.Bid) BidResponse {
	return BidResponse{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}
}

func NewBidResponseList(bids []entity.Bid) []BidResponse {
	bidsResponseList := make([]BidResponse, 0, len(bids))

	for i := range bids {
		bidsResponseList = append(bidsResponseList, NewBidResponse(bids[i]))
	}

	return bidsResponseList
}
