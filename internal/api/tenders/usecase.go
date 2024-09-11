package tenders

import (
	"context"

	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/pkg/queryparams"
)

type ListOpts struct {
}

type Usecase interface {
	Create(ctx context.Context, request dtos.CreateTenderRequest) (dtos.TenderResponse, error)
	Edit(ctx context.Context, id string, request dtos.EditTenderRequest) (dtos.TenderResponse, error)
	EditStatus(ctx context.Context, id string, request dtos.EditTenderStatusRequest) (dtos.TenderResponse, error)
	Rollback(ctx context.Context, id string, request dtos.RollbackTenderRequest) (dtos.TenderResponse, error)
	GetAll(ctx context.Context, filter TenderFilter, pagination queryparams.Pagination) ([]dtos.TenderResponse, error)
	GetTenderStatus(ctx context.Context, id string, request dtos.TenderStatus) (dtos.TenderResponse, error)
	FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.TenderResponse, error)
}
