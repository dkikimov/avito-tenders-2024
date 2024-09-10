package entities

import "avito-tenders/internal/entity"

type TenderStatusResponse struct {
	Status entity.TenderStatus `json:"tenderStatus"`
}
