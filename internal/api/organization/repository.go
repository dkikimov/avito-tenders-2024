package organization

import "context"

type Repository interface {
	IsOrganizationResponsible(ctx context.Context, organizationID int, username string) (bool, error)
}
