package employee

import (
	"context"

	"avito-tenders/internal/entity"
)

type Repository interface {
	FindById(ctx context.Context, id string) (entity.Employee, error)
}
