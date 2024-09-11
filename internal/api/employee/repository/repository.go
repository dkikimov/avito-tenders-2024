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

func NewRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) *Repository {
	return &Repository{db: db, getter: getter}
}

func (r Repository) FindById(ctx context.Context, id string) (entity.Employee, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx, `
	select id, username, first_name, last_name, created_at, updated_at from employee
	where id = $1`, id)
	if row.Err() != nil {
		return entity.Employee{}, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
	}

	var emp entity.Employee
	if err := row.StructScan(&emp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Employee{}, apperror.Unauthorized(apperror.ErrUserDoesNotExist)
		}

		slog.Error("couldn't scan employee found by id", "error", err)
		return entity.Employee{}, apperror.BadRequest(apperror.ErrInternal)
	}

	return emp, nil
}
