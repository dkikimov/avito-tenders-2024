package models

import "avito-tenders/internal/entity"

type EditTenderStatus struct {
	TenderID string
	Status   entity.TenderStatus
	Username string
}
