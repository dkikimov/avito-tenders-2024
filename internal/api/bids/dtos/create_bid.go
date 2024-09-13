package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type CreateBidRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	TenderID    string            `json:"tenderId"`
	AuthorType  entity.AuthorType `json:"authorType"`
	AuthorID    string            `json:"authorId"`
}

func (r CreateBidRequest) ToEntity() entity.Bid {
	return entity.Bid{
		Name:        r.Name,
		Description: r.Description,
		TenderID:    r.TenderID,
		AuthorType:  r.AuthorType,
		AuthorID:    r.AuthorID,
	}
}

func (r CreateBidRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(0, 100)),
		validation.Field(&r.Description, validation.Required, validation.Length(0, 500)),
		validation.Field(&r.TenderID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.AuthorType, validation.Required, r.AuthorType.ValidationRule()),
		validation.Field(&r.AuthorID, validation.Required, validation.Length(1, 100)),
	)
}
