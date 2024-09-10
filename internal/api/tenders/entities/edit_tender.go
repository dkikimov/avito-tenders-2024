package entities

import (
	"avito-tenders/internal/entity"
)

type EditTenderStatusRequest struct {
	Status   entity.TenderStatus `json:"status" valid:"required,tender_status"`
	Username string              `json:"username" valid:"required"`
}

type EditTender struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	ServiceType entity.ServiceType `json:"serviceType,omitempty" valid:"service_type"`
}

type EditTenderRequest struct {
	EditTender
	Username string `json:"username" valid:"required"`
}
