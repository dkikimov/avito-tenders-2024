package repository

import (
	"context"
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

func NewRepository(db *sqlx.DB, c *trmsqlx.CtxGetter) *Repository {
	return &Repository{db: db, getter: c}
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
		join employee e on e.id = b.author_id
		where e.username = $1
		limit $2 offset $3
`, req.Username, req.Limit, req.Offset)
	if err != nil {
		slog.Error("couldn't find bids by username", "error", err)
		return nil, apperror.InternalServerError(apperror.ErrInternal)
	}

	return bidsList, nil
}

func (r Repository) FindByID(ctx context.Context, id string) (entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) FindByTenderId(ctx context.Context, req models.FindByTenderId) ([]entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) Update(ctx context.Context, bid entity.Bid) (entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) FindByIDFromHistory(ctx context.Context, id string) (entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) SendFeedback(ctx context.Context, req models.SendFeedback) (entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}

func (r Repository) FindReviews(ctx context.Context, req models.FindReview) ([]entity.Bid, error) {
	// TODO implement me
	panic("implement me")
}
