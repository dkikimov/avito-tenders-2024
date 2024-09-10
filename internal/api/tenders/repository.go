package tenders

import (
	"context"

	"avito-tenders/internal/api/tenders/entities"
)

type Repository interface {
	Create(ctx context.Context, request entities.CreateTenderRequest) (entities.ResponseTender, error)
	Edit(ctx context.Context, id int, request entities.EditTenderRequest) (entities.ResponseTender, error)
	GetAll(ctx context.Context) ([]entities.ResponseTender, error)
	FindById(ctx context.Context, id int) (entities.ResponseTender, error)
	FindByUsername(ctx context.Context, username string) ([]entities.ResponseTender, error)
	EditStatus(ctx context.Context, id int, request entities.EditTenderStatusRequest) (entities.ResponseTender, error)
}
