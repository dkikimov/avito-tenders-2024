package models

import "avito-tenders/pkg/queryparams"

type FindByUsername struct {
	Username string
	queryparams.Pagination
}
