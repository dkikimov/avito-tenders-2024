package entities

import (
	"avito-tenders/internal/entity"
)

type CreateTenderRequest struct {
	Name            string              `json:"name" valid:"required"`
	Description     string              `json:"description" valid:"required"`
	ServiceType     string              `json:"serviceType" valid:"required"`
	Status          entity.TenderStatus `json:"status" valid:"required"`
	OrganizationId  int                 `json:"organizationId" valid:"required"`
	CreatorUsername string              `json:"creatorUsername" valid:"required"`
}
