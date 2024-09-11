package entities

import "avito-tenders/internal/entity"

type TenderStatus struct {
	Username string `json:"username"`
}

type TenderStatusResponse struct {
	Status entity.TenderStatus `json:"tenderStatus"`
}
