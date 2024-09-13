package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/pkg/queryparams"
)

type FindByTenderIDRequest struct {
	TenderID string `json:"tender_id"`
	Username string `json:"username"`
	queryparams.Pagination
}

func (r FindByTenderIDRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TenderID, validation.Required),
		validation.Field(&r.Username, validation.Required))
}
