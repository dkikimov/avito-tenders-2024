package fwcontext

import (
	"context"

	"avito-tenders/pkg/queryparams"
)

type CtxKey int

const (
	UsernameCtxKey CtxKey = iota
	PaginationCtxKey
)

func GetUsername(ctx context.Context) string {
	username, ok := ctx.Value(UsernameCtxKey).(string)
	if !ok {
		username = ""
	}

	return username
}

func GetPagination(ctx context.Context) queryparams.Pagination {
	pagination, ok := ctx.Value(PaginationCtxKey).(queryparams.Pagination)
	if !ok {
		pagination = queryparams.Pagination{
			Limit:  queryparams.DefaultLimit,
			Offset: queryparams.DefaultOffset,
		}
	}

	return pagination
}
