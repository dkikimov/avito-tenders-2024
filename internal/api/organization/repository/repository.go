package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"avito-tenders/internal/api/organization"
	"avito-tenders/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func (r repository) IsOrganizationResponsible(ctx context.Context, organizationID string, username string) (bool, error) {
	row := r.db.QueryRowxContext(ctx, `
		select o.id from organization_responsible o
		          join employee e on e.username = $1
		          where o.organization_id = $2 and o.user_id = e.id
`, username, organizationID)

	if err := row.Err(); err != nil {
		slog.Error("couldn't query is organization responsible", "error", err)
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	var id string
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		slog.Error("couldn't scan organization responsible", "error", err)
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	return true, nil
}

func NewRepository(db *sqlx.DB) organization.Repository {
	return &repository{db: db}
}
