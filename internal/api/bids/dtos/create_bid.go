package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type CreateBidRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	TenderId    string            `json:"tenderId"`
	AuthorType  entity.AuthorType `json:"authorType"`
	AuthorId    string            `json:"authorId"`
}

func (r CreateBidRequest) ToEntity() entity.Bid {
	return entity.Bid{
		Name:        r.Name,
		Description: r.Description,
		TenderId:    r.TenderId,
		AuthorType:  r.AuthorType,
		AuthorId:    r.AuthorId,
	}
}

func (r CreateBidRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Required, validation.Length(1, 500)),
		validation.Field(&r.TenderId, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.AuthorType, validation.Required, r.AuthorType.ValidationRule()),
		validation.Field(&r.AuthorId, validation.Required, validation.Length(1, 100)),
	)
}
