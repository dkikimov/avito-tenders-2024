package models

import "avito-tenders/internal/entity"

type EditTender struct {
	TenderID    string
	Name        string
	Description string
	ServiceType entity.ServiceType
	Username    string
}
