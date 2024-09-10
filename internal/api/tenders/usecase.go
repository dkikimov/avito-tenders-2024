package tenders

import (
	"context"

	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/query"
)

type ListOpts struct {
}

type Usecase interface {
	Create(ctx context.Context, request entities.CreateTenderRequest) (entities.ResponseTender, error)
	Edit(ctx context.Context, id string, request entities.EditTenderRequest) (entities.ResponseTender, error)
	FindList(ctx context.Context, id string) ([]entities.ResponseTender, error)
	FindById(ctx context.Context, id string) (entities.ResponseTender, error)
	FindByUsername(ctx context.Context, username string, pagination query.Pagination) ([]entities.ResponseTender, error)
}
