package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/entities"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/query"
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

		slog.Error("failed to insert tender", "error", row.Err())
		return entities.ResponseTender{}, apperror.InternalServerError(row.Err())
	}

	var result entities.ResponseTender
	if err := row.StructScan(&result); err != nil {
		return entities.ResponseTender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return result, nil
}

func (r repository) FindByUsername(ctx context.Context, username string, pagination query.Pagination) ([]entities.ResponseTender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	exist, err := r.DoesUserExist(ctx, tx, username)
	if err != nil || !exist {
		return nil, err
	}

	// Find user's tenders.
	var tenderList = make([]entities.ResponseTender, 0)
	err = tx.SelectContext(ctx, &tenderList, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders 
		where creator_username = $1
		limit $2
		offset $3`,
		username,
		pagination.Limit,
		pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tenderList, nil
}

func (r repository) FindById(ctx context.Context, id int) (entities.ResponseTender, error) {
	row := r.db.QueryRowxContext(ctx, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders 
		where id = $1`,
		id)

	if row.Err() != nil {
		slog.Error("failed to select", "error", row.Err())
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entities.ResponseTender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.BadRequest(apperror.ErrNotFound)
		}

		slog.Error("failed to scan", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r repository) GetAll(ctx context.Context) ([]entities.ResponseTender, error) {
	var tenderList = make([]entities.ResponseTender, 0)

	err := r.db.SelectContext(ctx, tenderList, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders`)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return tenderList, nil
}

func (r repository) EditStatus(ctx context.Context, id int, request entities.EditTenderStatusRequest) (entities.ResponseTender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	// Check does user exist.
	exist, err := r.DoesUserExist(ctx, tx, request.Username)
	if err != nil || !exist {
		return entities.ResponseTender{}, err
	}

	// Check does tender exist.
	row := r.db.QueryRowxContext(ctx, `select id from tenders where id = $1`, id)
	if row.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to select id tender", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tenderId int
	if err := row.Scan(&tenderId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to scan tender id", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	// Try to update tender.
	// If returns 0, that means creator_username doesn't have enough permissions.
	row = r.db.QueryRowxContext(ctx, `
		update tenders set status = $1, version = version + 1 where id = $2 and creator_username = $3
		returning id, name, description, service_type, status, organization_id, version, created_at 
`, request.Status, id, request.Username)
	if row.Err() != nil {
		slog.Error("failed to update status", "error", row.Err())
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entities.ResponseTender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.Forbidden(errors.New("user doesn't have enough permissions"))
		}

		slog.Error("failed to scan", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r repository) Edit(ctx context.Context, id int, request entities.EditTenderRequest) (entities.ResponseTender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	// Check does user exist.
	exist, err := r.DoesUserExist(ctx, tx, request.Username)
	if err != nil || !exist {
		return entities.ResponseTender{}, err
	}

	// Check does tender exist.
	row := r.db.QueryRowxContext(ctx, `select id from tenders where id = $1`, id)
	if row.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to select id tender", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tenderId int
	if err := row.Scan(&tenderId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to scan tender id", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	// Try to update tender.
	// If returns 0, that means creator_username doesn't have enough permissions.
	row = r.db.QueryRowxContext(
		ctx,
		`
		update tenders set 
		                   name = CASE WHEN $1 != '' THEN $1 ELSE name END,
		                   description = CASE WHEN $2 != '' THEN $2 ELSE description END,
		                   service_type = CASE WHEN $3 != '' THEN $3 ELSE service_type END,
		                   version = version + 1 
		               		where id = $4 and creator_username = $5
		returning id, name, description, service_type, status, organization_id, version, created_at 
`,
		request.Name,
		request.Description,
		request.ServiceType,
		id,
		request.Username,
	)
	if row.Err() != nil {
		slog.Error("failed to update status", "error", row.Err())
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entities.ResponseTender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.Forbidden(errors.New("user doesn't have enough permissions"))
		}

		slog.Error("failed to scan", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r repository) DoesUserExist(ctx context.Context, tx *sqlx.Tx, username string) (bool, error) {
	var id int
	row := tx.QueryRowxContext(ctx, "select id from employee e where username = $1", username)
	if row.Err() != nil {
		slog.Error("failed to select", row.Err())
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
		}

		slog.Error("failed to scan", row.Err())
		return false, apperror.InternalServerError(apperror.ErrInternal)
	}

	return true, nil
}

func (r repository) Rollback(ctx context.Context, id int, request entities.RollbackTender) (entities.ResponseTender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	// Check does user exist.
	exist, err := r.DoesUserExist(ctx, tx, request.Username)
	if err != nil || !exist {
		return entities.ResponseTender{}, err
	}

	// Check does tender exist.
	row := r.db.QueryRowxContext(ctx, `select id from tenders where id = $1`, id)
	if row.Err() != nil {
		slog.Error("failed to select id tender", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tenderId int
	if err := row.Scan(&tenderId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to scan tender id", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	// Get old version
	row = r.db.QueryRowxContext(ctx, `
		select id, name, description, service_type, status, organization_id, version, created_at from tender_history 
		where tender_id = $1 and version = $2`,
		id, request.Version)
	if row.Err() != nil {
		slog.Error("failed to select id tender", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var oldTender entities.ResponseTender
	if err := row.StructScan(&oldTender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to scan old tender", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	// Try to update tender.
	// If returns 0, that means creator_username doesn't have enough permissions.
	row = r.db.QueryRowxContext(
		ctx,
		`
		update tenders set 
		                   name = $1,
		                   description = $2,
		                   service_type = $3,
		                   status = $4, 
		                   organization_id = $5, 
		                   created_at = $6,
		                   version = version + 1 
		               		where id = $7 and creator_username = $8
		returning id, name, description, service_type, status, organization_id, version, created_at 
`,
		oldTender.Name,
		oldTender.Description,
		oldTender.ServiceType,
		oldTender.Status,
		oldTender.OrganizationId,
		oldTender.CreatedAt,
		id,
		request.Username,
	)
	if row.Err() != nil {
		slog.Error("failed to update status", "error", row.Err())
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entities.ResponseTender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ResponseTender{}, apperror.Forbidden(errors.New("user doesn't have enough permissions"))
		}

		slog.Error("failed to scan", "error", err)
		return entities.ResponseTender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}
