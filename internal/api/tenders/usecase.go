package tenders

import (
	"context"

	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/queryparams"
)

type ListOpts struct {
}

type Usecase interface {
	Create(ctx context.Context, request entities.CreateTenderRequest) (entities.ResponseTender, error)
	Edit(ctx context.Context, id string, request entities.EditTenderRequest) (entities.ResponseTender, error)
	EditStatus(ctx context.Context, id string, request entities.EditTenderStatusRequest) (entities.ResponseTender, error)
	Rollback(ctx context.Context, id string, request entities.RollbackTenderRequest) (entities.ResponseTender, error)
	GetAll(ctx context.Context, filter TenderFilter, pagination queryparams.Pagination) ([]entities.ResponseTender, error)
	GetTenderStatus(ctx context.Context, id string, request entities.TenderStatus) (entities.ResponseTender, error)
	FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]entities.ResponseTender, error)
}
