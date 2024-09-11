package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/pkg/queryparams"
)

type FindReviewsRequest struct {
	TenderId          string `json:"tenderId"`
	AuthorUsername    string `json:"authorUsername"`
	RequesterUsername string `json:"requesterUsername"`
	queryparams.Pagination
}

func (r FindReviewsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TenderId, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.AuthorUsername, validation.Required),
		validation.Field(&r.RequesterUsername, validation.Required))
}
