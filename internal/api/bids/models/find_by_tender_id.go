package models

import (
	"avito-tenders/pkg/queryparams"
)

type FindByTenderId struct {
	TenderId string
	queryparams.Pagination
}
