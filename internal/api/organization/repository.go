package organization

import (
	"context"

	"avito-tenders/internal/entity"
)

type Repository interface {
	// IsOrganizationResponsible checks if user is responsible in given organization.
	IsOrganizationResponsible(ctx context.Context, organizationID, username string) (bool, error)

	// GetUserOrganization returns user's organization.
	GetUserOrganization(ctx context.Context, userID string) (entity.Organization, error)

	// GetOrganizationResponsible returns slice of responsible ids.
	GetOrganizationResponsible(ctx context.Context, organizationID string) ([]string, error)

	// FindByID returns organization found by organization id.
	FindByID(ctx context.Context, organizationID string) (entity.Organization, error)
}
