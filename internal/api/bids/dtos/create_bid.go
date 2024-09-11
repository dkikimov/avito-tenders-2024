package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type CreateBidRequest struct {
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Status          entity.BidStatus `json:"status"`
	TenderId        string           `json:"tender_id"`
	OrganizationId  int              `json:"organization_id"`
	CreatorUsername string           `json:"creator_username"`
}

func (r CreateBidRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Status, validation.Required, validation.In(
			entity.BidApproved,
			entity.BidCanceled,
			entity.BidCreated,
			entity.BidPublished,
			entity.BidRejected,
		)),
		validation.Field(&r.OrganizationId, validation.Required),
		validation.Field(&r.CreatorUsername, validation.Required),
	)
}
