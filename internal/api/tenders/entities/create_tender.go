package entities

import (
	"avito-tenders/internal/entity"
)

type CreateTenderRequest struct {
	Name            string              `json:"name" valid:"required"`
	Description     string              `json:"description" valid:"required"`
	ServiceType     entity.ServiceType  `json:"serviceType" valid:"required,service_type"`
	Status          entity.TenderStatus `json:"status" valid:"required,tender_status"`
	OrganizationId  int                 `json:"organizationId" valid:"required"`
	CreatorUsername string              `json:"creatorUsername" valid:"required"`
}
