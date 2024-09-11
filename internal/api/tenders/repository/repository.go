package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"avito-tenders/internal/api/tenders"
	"avito-tenders/internal/api/tenders/models"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
	"avito-tenders/pkg/queryparams"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) FindByIDFromHistory(ctx context.Context, id string, version int) (entity.Tender, error) {
	row := r.db.QueryRowxContext(ctx, `
		select tender_id as id, name, description, service_type, status, organization_id, version, created_at from tenders_history
		where tender_id = $1 and version = $2`,
		id, version)
	if row.Err() != nil {
		return entity.Tender{}, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var oldTender entity.Tender
	if err := row.StructScan(&oldTender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Tender{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("failed to scan old tender", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return oldTender, nil
}

func (r Repository) Create(ctx context.Context, tender entity.Tender) (entity.Tender, error) {
	row := r.db.QueryRowxContext(ctx, `
		INSERT INTO tenders(name, description, service_type, status, organization_id, creator_username) 
		VALUES($1,$2,$3,$4,$5,$6)
		returning id, name, description, service_type, status, organization_id, version, created_at`,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status.String(),
		tender.OrganizationId,
		tender.CreatorUsername)
	if row.Err() != nil {
		var pgError *pgconn.PgError
		if errors.As(row.Err(), &pgError) {
			if pgError.Code == "23503" {
				return entity.Tender{}, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
			}
		}

		slog.Error("failed to insert tender", "error", row.Err())
		return entity.Tender{}, apperror.InternalServerError(row.Err())
	}

	var result entity.Tender
	if err := row.StructScan(&result); err != nil {
		return entity.Tender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return result, nil
}

func (r Repository) Update(ctx context.Context, tender entity.Tender) (entity.Tender, error) {
	row := r.db.QueryRowxContext(ctx, `
		update tenders set
		                   name = $1,
		                   description = $2,
		                   service_type = $3,
		                   status = $4,
		                   organization_id = $5,
		                   version = version + 1
		               where id = $6
		returning id, name, description, service_type, status, organization_id, creator_username, created_at, version`,
		tender.Name,
		tender.Description,
		tender.ServiceType,
		tender.Status.String(),
		tender.OrganizationId,
		tender.Id)
	if row.Err() != nil {
		var pgError *pgconn.PgError
		if errors.As(row.Err(), &pgError) {
			if errors.Is(row.Err(), sql.ErrNoRows) {
				return entity.Tender{}, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
			}
		}

		slog.Error("failed to update tender", "error", row.Err())
		return entity.Tender{}, apperror.InternalServerError(row.Err())
	}

	var result entity.Tender
	if err := row.StructScan(&result); err != nil {
		return entity.Tender{}, fmt.Errorf("failed to scan: %w", err)
	}

	return result, nil
}

func (r Repository) FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]entity.Tender, error) {
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
	var tenderList = make([]entity.Tender, 0)
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

func (r Repository) FindById(ctx context.Context, id string) (entity.Tender, error) {
	row := r.db.QueryRowxContext(ctx, `
		select id, name, description, service_type, status, organization_id, version, created_at from tenders 
		where id = $1`,
		id)

	if row.Err() != nil {
		slog.Error("failed to select", "error", row.Err())
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entity.Tender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Tender{}, apperror.BadRequest(apperror.ErrNotFound)
		}

		slog.Error("failed to scan", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r Repository) GetAll(ctx context.Context, filter tenders.TenderFilter, pagination queryparams.Pagination) ([]entity.Tender, error) {
	var filterValues = make([]interface{}, 0)

	query := strings.Builder{}
	query.WriteString(`select id, name, description, service_type, status, organization_id, version, created_at from tenders 
    					where status = 'Published' `)

	if len(filter.ServiceTypes) > 0 {
		query.WriteString("and service_type IN (")

		for i, service := range filter.ServiceTypes {
			query.WriteString(fmt.Sprintf("$%d", len(filterValues)+1))
			filterValues = append(filterValues, service)
			if i != len(filter.ServiceTypes)-1 {
				query.WriteString(",")
			}
		}
		query.WriteString(") ")
	}

	query.WriteString(fmt.Sprintf("limit $%d offset $%d", len(filterValues)+1, len(filterValues)+2))
	filterValues = append(filterValues, pagination.Limit, pagination.Offset)

	var tenderList = make([]entity.Tender, 0)
	err := r.db.SelectContext(ctx, &tenderList, query.String(), filterValues...)
	if err != nil {
		slog.Error("failed to get all tenders", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tenderList, nil
}

func (r Repository) EditStatus(ctx context.Context, request models.EditTenderStatus) (entity.Tender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	// Check does user exist.
	exist, err := r.DoesUserExist(ctx, tx, request.Username)
	if err != nil || !exist {
		return entity.Tender{}, err
	}

	// Check does tender exist.
	if _, err := r.FindById(ctx, request.TenderId); err != nil {
		return entity.Tender{}, err
	}

	// Try to update tender.
	// If returns 0, that means creator_username doesn't have enough permissions.
	row := r.db.QueryRowxContext(ctx, `
		update tenders set status = $1, version = version + 1 where id = $2 and creator_username = $3
		returning id, name, description, service_type, status, organization_id, version, created_at 
`, request.Status, request.TenderId, request.Username)
	if row.Err() != nil {
		slog.Error("failed to update status", "error", row.Err())
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entity.Tender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Tender{}, apperror.Forbidden(errors.New("user doesn't have enough permissions"))
		}

		slog.Error("failed to scan", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r Repository) Edit(ctx context.Context, editTender models.EditTender) (entity.Tender, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	// Check does user exist.
	exist, err := r.DoesUserExist(ctx, tx, editTender.Username)
	if err != nil || !exist {
		return entity.Tender{}, err
	}

	// Check does tender exist.
	if _, err := r.FindById(ctx, editTender.TenderID); err != nil {
		return entity.Tender{}, err
	}

	// Try to update tender.
	// If returns 0, that means creator_username doesn't have enough permissions.
	row := r.db.QueryRowxContext(
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
		editTender.Name,
		editTender.Description,
		editTender.ServiceType,
		editTender.TenderID,
		editTender.Username,
	)
	if row.Err() != nil {
		slog.Error("failed to update status", "error", row.Err())
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	var tender entity.Tender
	if err := row.StructScan(&tender); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Tender{}, apperror.Forbidden(errors.New("user doesn't have enough permissions"))
		}

		slog.Error("failed to scan", "error", err)
		return entity.Tender{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return tender, nil
}

func (r Repository) DoesUserExist(ctx context.Context, tx *sqlx.Tx, username string) (bool, error) {
	var id string
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
