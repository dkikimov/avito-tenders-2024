package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type UpdateStatusRequest struct {
	BidID    string           `json:"bidId"`
	Status   entity.BidStatus `json:"status"`
	Username string           `json:"username"`
}

func (r UpdateStatusRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Status, validation.Required, validation.In(entity.BidCreated, entity.BidPublished, entity.BidCanceled)),
		validation.Field(&r.Username, validation.Required))
}
