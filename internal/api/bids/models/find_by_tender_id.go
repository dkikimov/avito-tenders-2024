package models

import (
	"avito-tenders/pkg/queryparams"
)

type FindByTenderId struct {
	TenderId string
	Username string
	queryparams.Pagination
}
