package tenders

import (
	"context"

	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/queryparams"
)

type TenderFilter struct {
	ServiceTypes []entity.ServiceType
}

type Repository interface {
	Create(ctx context.Context, request dtos.CreateTenderRequest) (dtos.TenderResponse, error)
	Edit(ctx context.Context, id int, request dtos.EditTenderRequest) (dtos.TenderResponse, error)
	GetAll(ctx context.Context, filter TenderFilter, pagination queryparams.Pagination) ([]dtos.TenderResponse, error)
	FindById(ctx context.Context, id int) (dtos.TenderResponse, error)
	FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.TenderResponse, error)
	EditStatus(ctx context.Context, id int, request dtos.EditTenderStatusRequest) (dtos.TenderResponse, error)
	Rollback(ctx context.Context, id int, request dtos.RollbackTender) (dtos.TenderResponse, error)
}
