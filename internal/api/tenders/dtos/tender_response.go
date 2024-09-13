package dtos

import (
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/types"
)

type TenderResponse struct {
	Id             string              `json:"id" db:"id"`
	Name           string              `json:"name" db:"name"`
	Description    string              `json:"description" db:"description"`
	ServiceType    entity.ServiceType  `json:"serviceType" db:"service_type"`
	Status         entity.TenderStatus `json:"status" db:"status"`
	OrganizationId string              `json:"organizationId" db:"organization_id"`
	Version        int                 `json:"version" db:"version"`
	CreatedAt      types.RFC3339Time   `json:"createdAt" db:"created_at"`
}

func NewTenderResponse(tender entity.Tender) TenderResponse {
	return TenderResponse{
		Id:             tender.Id,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.Status,
		OrganizationId: tender.OrganizationId,
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
