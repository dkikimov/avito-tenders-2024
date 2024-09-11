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
