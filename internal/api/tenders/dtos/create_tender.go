package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type CreateTenderRequest struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	ServiceType     entity.ServiceType  `json:"serviceType"`
	Status          entity.TenderStatus `json:"status"`
	OrganizationID  string              `json:"organizationId"`
	CreatorUsername string              `json:"creatorUsername"`
}

func (c CreateTenderRequest) ToEntity() entity.Tender {
	return entity.Tender{
		Name:            c.Name,
		Description:     c.Description,
		ServiceType:     c.ServiceType,
		Status:          c.Status,
		OrganizationID:  c.OrganizationID,
		CreatorUsername: c.CreatorUsername,
	}
}

func (c CreateTenderRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&c.Description, validation.Required),
		validation.Field(&c.ServiceType, validation.Required, c.ServiceType.ValidationRule()),
		validation.Field(&c.Status, validation.Required, c.Status.ValidationRule()),
		validation.Field(&c.OrganizationID, validation.Required),
		validation.Field(&c.CreatorUsername, validation.Required),
	)
}
