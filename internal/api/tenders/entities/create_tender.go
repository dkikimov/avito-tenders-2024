package entities

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type CreateTenderRequest struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	ServiceType     entity.ServiceType  `json:"serviceType"`
	Status          entity.TenderStatus `json:"status"`
	OrganizationId  int                 `json:"organizationId"`
	CreatorUsername string              `json:"creatorUsername"`
}

func (c CreateTenderRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&c.Description, validation.Required),
		validation.Field(&c.ServiceType, validation.Required, c.ServiceType.ValidationRules()),
		validation.Field(&c.Status, validation.Required, c.Status.ValidationRules()),
		validation.Field(&c.OrganizationId, validation.Required),
		validation.Field(&c.CreatorUsername, validation.Required),
	)
}
