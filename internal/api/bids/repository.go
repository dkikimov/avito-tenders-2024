package bids

import (
	"context"

	"avito-tenders/internal/api/bids/models"
	"avito-tenders/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, bid entity.Bid) (entity.Bid, error)
	FindByUsername(ctx context.Context, req models.FindByUsername) ([]entity.Bid, error)
	FindByID(ctx context.Context, id string) (entity.Bid, error)
	FindByTenderId(ctx context.Context, req models.FindByTenderId) ([]entity.Bid, error)
	Update(ctx context.Context, bid entity.Bid) (entity.Bid, error)
	FindByIDFromHistory(ctx context.Context, id string) (entity.Bid, error)
	SendFeedback(ctx context.Context, req models.SendFeedback) (entity.Bid, error)
	FindReviews(ctx context.Context, req models.FindReview) ([]entity.Bid, error)
}
