package employee

import (
	"context"

	"avito-tenders/internal/entity"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (entity.Employee, error)
	FindByID(ctx context.Context, id string) (entity.Employee, error)
}
