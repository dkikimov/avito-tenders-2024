package models

import "avito-tenders/internal/entity"

type EditTenderStatus struct {
	TenderId string
	Status   entity.TenderStatus
	Username string
}
