package models

import (
	"avito-tenders/pkg/queryparams"
)

type FindByTenderID struct {
	TenderID string
	queryparams.Pagination
}
