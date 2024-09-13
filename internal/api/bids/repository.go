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
	FindByIDFromHistory(ctx context.Context, id string, version int) (entity.Bid, error)
	SendFeedback(ctx context.Context, req models.SendFeedback) error
	FindReviews(ctx context.Context, req models.FindReview) ([]entity.Review, error)
	SubmitApproveDecision(ctx context.Context, bidId string, userId string) error
	GetBidApproveAmount(ctx context.Context, bidId string) (int, error)
	FindBidsByOrganization(ctx context.Context, organizationId string) ([]entity.Bid, error)
}
