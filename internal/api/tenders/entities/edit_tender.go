package entities

import (
	"avito-tenders/internal/entity"
)

type EditTenderStatusRequest struct {
	Status   entity.TenderStatus `json:"status" valid:"required,tender_status"`
	Username string              `json:"username" valid:"required"`
}

type EditTenderRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ServiceType *string `json:"serviceType,omitempty"`
}
