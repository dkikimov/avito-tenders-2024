package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) tenders.Repository {
	return &repository{db: db}
}

func (r repository) Create(ctx context.Context, request entities.CreateTenderRequest) (entities.ResponseTender, error) {
	row := r.db.QueryRowxContext(ctx, `
		INSERT INTO tenders(name, description, service_type, status, organization_id, creator_username) 
		VALUES($1,$2,$3,$4,$5,$6)
		returning id, name, description, service_type, status, organization_id, version, created_at`,
		request.Name,
		request.Description,
		request.ServiceType,
		request.Status.String(),
		request.OrganizationId,
		request.CreatorUsername)
	if row.Err() != nil {
		var pgError *pgconn.PgError
		if errors.As(row.Err(), &pgError) {
			if pgError.Code == "23503" {
				return entities.ResponseTender{}, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
			}
		}

		return entities.ResponseTender{}, apperror.InternalServerError(row.Err())
	}

	var result entities.ResponseTender
	if err := row.StructScan(&result); err != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return result, nil
}

func (r repository) FindByUsername(ctx context.Context, username string) ([]entities.ResponseTender, error) {
	var tenders = make([]entities.ResponseTender, 0)

	err := r.db.SelectContext(ctx, tenders, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders 
		where creator_username = ?`,
		username)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return tenders, nil
}

func (r repository) FindById(ctx context.Context, id int) (entities.ResponseTender, error) {
	var tenders entities.ResponseTender
	err := r.db.SelectContext(ctx, tenders, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders 
		where id = ?`,
		id)
	if err != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to select: %w", err)
	}

	return tenders, nil
}

func (r repository) GetAll(ctx context.Context) ([]entities.ResponseTender, error) {
	var tenders = make([]entities.ResponseTender, 0)

	err := r.db.SelectContext(ctx, tenders, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders`)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return tenders, nil
}

func (r repository) EditStatus(ctx context.Context, id int, request entities.EditTenderStatusRequest) (entities.ResponseTender, error) {
	row := r.db.QueryRowxContext(ctx, `
		update tenders set status = ? where id = ? and creator_username = ?
		returning id, name, description, service_type, status, organization_id, version, created_at 
`)
	if row.Err() != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to select: %w", row.Err())
	}

	var tender entities.ResponseTender
	if err := row.Scan(&tender); err != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return tender, nil
}

func (r repository) Edit(ctx context.Context, id int, request entities.EditTenderRequest) (entities.ResponseTender, error) {
	row := r.db.QueryRowxContext(ctx, `
		update tenders set status = ? where id = ? and creator_username = ?
		returning id, name, description, service_type, status, organization_id, version, created_at 
`)
	if row.Err() != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to select: %w", row.Err())
	}

	var tender entities.ResponseTender
	if err := row.Scan(&tender); err != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return tender, nil
}
