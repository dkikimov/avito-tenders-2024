package organization

import (
	"context"

	"avito-tenders/internal/entity"
)

type Repository interface {
	// IsOrganizationResponsible checks if user is responsible in given organization.
	IsOrganizationResponsible(ctx context.Context, organizationID string, username string) (bool, error)

	// GetUserOrganization returns user's organization.
	GetUserOrganization(ctx context.Context, userId string) (entity.Organization, error)

	// GetOrganizationResponsible returns slice of responsible ids.
	GetOrganizationResponsible(ctx context.Context, organizationID string) ([]string, error)

	// FindById returns organization found by organization id.
	FindById(ctx context.Context, organizationID string) (entity.Organization, error)
}
