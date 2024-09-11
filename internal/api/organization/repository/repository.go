package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
)

type Repository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewRepository(db *sqlx.DB, g *trmsqlx.CtxGetter) *Repository {
	return &Repository{db: db, getter: g}
}

func (r Repository) GetUserOrganization(ctx context.Context, userId string) (entity.Organization, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx, `
		select o.id, o.name, o.description, o.type, o.created_at, o.updated_at from organization_responsible r
		                join organization o on o.id = r.organization_id
		                where user_id = $1`, userId)
	if row.Err() != nil {
		return entity.Organization{}, apperror.Unauthorized(apperror.ErrUnauthorized)
	}

	var org entity.Organization
	if err := row.StructScan(&org); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Organization{}, apperror.Unauthorized(apperror.ErrUnauthorized)
		}

		slog.Error("couldn't scan organization found by user id", "error", err)
		return entity.Organization{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return org, nil
}

func (r Repository) IsOrganizationResponsible(ctx context.Context, organizationID string, username string) (bool, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx, `
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

func (r Repository) FindById(ctx context.Context, organizationID string) (entity.Organization, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx, `
		select id, name, description, type, created_at, updated_at from organization
		where id = $1`, organizationID)
	if row.Err() != nil {
		return entity.Organization{}, apperror.Unauthorized(apperror.ErrOrganizationDoesNotExist)
	}

	var org entity.Organization
	if err := row.StructScan(&org); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Organization{}, apperror.Unauthorized(apperror.ErrOrganizationDoesNotExist)
		}

		slog.Error("couldn't scan organization found by id", "error", err)
		return entity.Organization{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return org, nil
}
