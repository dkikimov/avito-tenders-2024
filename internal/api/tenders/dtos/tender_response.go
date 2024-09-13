package dtos

import (
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/types"
)

type TenderResponse struct {
	ID             string              `json:"id" db:"id"`
	Name           string              `json:"name" db:"name"`
	Description    string              `json:"description" db:"description"`
	ServiceType    entity.ServiceType  `json:"serviceType" db:"service_type"`
	Status         entity.TenderStatus `json:"status" db:"status"`
	OrganizationID string              `json:"organizationId" db:"organization_id"`
	Version        int                 `json:"version" db:"version"`
	CreatedAt      types.RFC3339Time   `json:"createdAt" db:"created_at"`
}

func NewTenderResponse(tender entity.Tender) TenderResponse {
	return TenderResponse{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.Status,
		OrganizationID: tender.OrganizationID,
		Version:        tender.Version,
		CreatedAt:      types.RFCFromTime(tender.CreatedAt),
	}
}

func NewTenderResponseList(tendersList []entity.Tender) []TenderResponse {
	dtoTenders := make([]TenderResponse, 0, len(tendersList))
	for i := range tendersList {
		dtoTenders = append(dtoTenders, NewTenderResponse(tendersList[i]))
	}

	return dtoTenders
}
