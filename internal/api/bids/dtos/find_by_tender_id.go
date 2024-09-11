package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/pkg/queryparams"
)

type FindByTenderIdRequest struct {
	TenderId string `json:"tender_id"`
	Username string `json:"username"`
	queryparams.Pagination
}

func (r FindByTenderIdRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TenderId, validation.Required),
		validation.Field(&r.Username, validation.Required))
}
