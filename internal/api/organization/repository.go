package organization

import "context"

type Repository interface {
	IsOrganizationResponsible(ctx context.Context, organizationID string, username string) (bool, error)
}
