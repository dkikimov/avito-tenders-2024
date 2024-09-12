package dtos

import "avito-tenders/internal/entity"

type BidStatusResponse struct {
	BidStatus entity.BidStatus `json:"bidStatus"`
}
