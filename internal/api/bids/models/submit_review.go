package models

import "avito-tenders/pkg/queryparams"

type FindReview struct {
	TenderId          string `json:"tenderId"`
	AuthorUsername    string `json:"authorUsername"`
	RequesterUsername string `json:"requesterUsername"`
	queryparams.Pagination
}
