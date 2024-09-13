package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

	"avito-tenders/internal/api/bids/models"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/apperror"
)

type Repository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func (r Repository) GetBidApproveAmount(ctx context.Context, bidId string) (int, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx,
		`select count(bid_id) from bids_approvals where bid_id = $1`, bidId)
	if err := row.Err(); err != nil {
		return 0, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	return count, nil
}

func NewRepository(db *sqlx.DB, c *trmsqlx.CtxGetter) *Repository {
	return &Repository{db: db, getter: c}
}

func (r Repository) SubmitApproveDecision(ctx context.Context, bidId string, userId string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("couldn't start transaction", "error", err)
		return apperror.InternalServerError(apperror.ErrInternal)
	}
	defer tx.Rollback()

	tr := r.getter.DefaultTrOrDB(ctx, tx)

	_, err = tr.ExecContext(ctx,
		`insert into bids_approvals (bid_id, user_id)
				values ($1, $2) 
				on conflict do nothing`, bidId, userId)
	if err != nil {
		return apperror.BadRequest(apperror.ErrInvalidInput)
	}

	return nil
}

func (r Repository) Create(ctx context.Context, bid entity.Bid) (entity.Bid, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.db)

	row := tr.QueryRowxContext(
		ctx,
		`
		insert into bids(name, description, status, tender_id, author_type, author_id) 
		VALUES ($1, $2, $3, $4, $5, $6)
		returning id, name, description, status, tender_id, author_type, author_id, version, created_at
`,
		bid.Name,
		bid.Description,
		entity.BidCreated,
		bid.TenderId,
		bid.AuthorType,
		bid.AuthorId,
	)

	if row.Err() != nil {
		return entity.Bid{}, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var createdBid entity.Bid
	if err := row.StructScan(&createdBid); err != nil {
		slog.Error("couldn't scan created bid", "error", err)
		return entity.Bid{}, apperror.InternalServerError(err)
	}

	return createdBid, nil
}

func (r Repository) FindByUsername(ctx context.Context, req models.FindByUsername) ([]entity.Bid, error) {
	var bidsList []entity.Bid

	err := r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &bidsList, `
		select b.id, b.name, b.description, b.status, b.tender_id, b.author_type, b.author_id, b.version, b.created_at from bids b 
		join employee e on e.username = $1
		where author_id = e.id
		order by name
		limit $2 offset $3
`, req.Username, req.Limit, req.Offset)
	if err != nil {
		slog.Error("couldn't find bids by username", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}

	return bidsList, nil
}

func (r Repository) FindByID(ctx context.Context, id string) (entity.Bid, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx,
		`select id, name, description, status, tender_id, author_type, author_id, version, created_at from bids
				where id = $1`, id)
	if row.Err() != nil {
		return entity.Bid{}, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var foundBid entity.Bid
	if err := row.StructScan(&foundBid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Bid{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("couldn't scan found bid row", "error", err)
		return entity.Bid{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return foundBid, nil
}

func (r Repository) FindByTenderId(ctx context.Context, req models.FindByTenderId) ([]entity.Bid, error) {
	var bidsList []entity.Bid

	err := r.getter.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &bidsList, `
		select id, name, description, status, tender_id, author_type, author_id, version, created_at from bids
		where tender_id = $1
		order by name
		limit $2 offset $3
`, req.TenderId, req.Limit, req.Offset)
	if err != nil {
		slog.Error("couldn't find bids by tender id", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}

	return bidsList, nil
}

func (r Repository) Update(ctx context.Context, bid entity.Bid) (entity.Bid, error) {
	tr := r.getter.DefaultTrOrDB(ctx, r.db)

	row := tr.QueryRowxContext(
		ctx,
		`
		update bids set 
		                name = $1,
		                description = $2,
		                status = $3,
		                tender_id = $4,
		                author_type = $5, 
		                author_id = $6,
		                version = version + 1
		            where id = $7
		returning id, name, description, status, tender_id, author_type, author_id, version, created_at
`,
		bid.Name,
		bid.Description,
		bid.Status,
		bid.TenderId,
		bid.AuthorType,
		bid.AuthorId,
		bid.Id,
	)
	if row.Err() != nil {
		return entity.Bid{}, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var updatedBid entity.Bid
	if err := row.StructScan(&updatedBid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Bid{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("couldn't scan updated bid", "error", err)
		return entity.Bid{}, apperror.InternalServerError(err)
	}

	return updatedBid, nil
}

func (r Repository) FindByIDFromHistory(ctx context.Context, id string, version int) (entity.Bid, error) {
	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRowxContext(ctx, `
select bid_id as id, name, description, status, tender_id, author_type, author_id, version, created_at
		from bids_history
		where bid_id = $1 and version = $2`, id, version)
	if err := row.Err(); err != nil {
		return entity.Bid{}, apperror.BadRequest(apperror.ErrInvalidInput)
	}

	var foundBid entity.Bid
	if err := row.StructScan(&foundBid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Bid{}, apperror.NotFound(apperror.ErrNotFound)
		}

		slog.Error("couldn't scan found bid row from history", "error", err)
		return entity.Bid{}, apperror.InternalServerError(apperror.ErrInternal)
	}

	return foundBid, nil
}

func (r Repository) SendFeedback(ctx context.Context, req models.SendFeedback) (entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) FindReviews(ctx context.Context, req models.FindReview) ([]entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}
