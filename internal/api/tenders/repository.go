package tenders

import (
	"context"

	"avito-tenders/internal/entity"
	"avito-tenders/pkg/queryparams"
)

type TenderFilter struct {
	ServiceTypes []entity.ServiceType
}

type Repository interface {
	Create(ctx context.Context, tender entity.Tender) (entity.Tender, error)
	Update(ctx context.Context, tender entity.Tender) (entity.Tender, error)
	GetAll(ctx context.Context, filter TenderFilter, pagination queryparams.Pagination) ([]entity.Tender, error)
	FindById(ctx context.Context, id string) (entity.Tender, error)
	FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]entity.Tender, error)
	FindByIDFromHistory(ctx context.Context, id string, version int) (entity.Tender, error)
}
