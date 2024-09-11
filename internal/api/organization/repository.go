package organization

import (
	"context"

	"avito-tenders/internal/entity"
)

type Repository interface {
	IsOrganizationResponsible(ctx context.Context, organizationID string, username string) (bool, error)
	GetUserOrganization(ctx context.Context, userId string) (entity.Organization, error)
	FindById(ctx context.Context, organizationID string) (entity.Organization, error)
}
