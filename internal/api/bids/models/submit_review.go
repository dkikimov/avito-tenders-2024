package models

import (
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/queryparams"
)

type FindReview struct {
	Bids []entity.Bid
	queryparams.Pagination
}
