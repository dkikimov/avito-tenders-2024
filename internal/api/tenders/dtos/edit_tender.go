package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type EditTenderStatusRequest struct {
	Status   entity.TenderStatus `json:"status" valid:"required,tender_status"`
	Username string              `json:"username" valid:"required"`
}

func (r EditTenderStatusRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Status, validation.Required, r.Status.ValidationRule()),
		validation.Field(&r.Username, validation.Required),
	)
}

type EditTender struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	ServiceType entity.ServiceType `json:"serviceType,omitempty"`
}

func (t EditTender) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.ServiceType, t.ServiceType.ValidationRule()))
}

type EditTenderRequest struct {
	EditTender
	Username string `json:"username" valid:"required"`
}
