package dtos

import (
	"time"

	"avito-tenders/internal/entity"
)

type TenderResponse struct {
	Id             int                 `db:"id"`
	Name           string              `db:"name"`
	Description    string              `db:"description"`
	ServiceType    entity.ServiceType  `db:"service_type"`
	Status         entity.TenderStatus `db:"status"`
	OrganizationId int                 `db:"organization_id"`
	Version        int                 `db:"version"`
	CreatedAt      time.Time           `db:"created_at"`
}
