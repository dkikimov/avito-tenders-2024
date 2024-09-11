package dtos

import (
	"time"

	"avito-tenders/internal/entity"
)

type TenderResponse struct {
	Id             string              `db:"id"`
	Name           string              `db:"name"`
	Description    string              `db:"description"`
	ServiceType    entity.ServiceType  `db:"service_type"`
	Status         entity.TenderStatus `db:"status"`
	OrganizationId string              `db:"organization_id"`
	Version        int                 `db:"version"`
	CreatedAt      time.Time           `db:"created_at"`
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
		CreatedAt:      tender.CreatedAt,
	}
}

func NewTenderResponseList(tendersList []entity.Tender) []TenderResponse {
	dtoTenders := make([]TenderResponse, 0, len(tendersList))
	for i := range tendersList {
		dtoTenders = append(dtoTenders, NewTenderResponse(tendersList[i]))
	}

	return dtoTenders
}
