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
	FindByTenderID(ctx context.Context, req models.FindByTenderID) ([]entity.Bid, error)
	Update(ctx context.Context, bid entity.Bid) (entity.Bid, error)
	FindByIDFromHistory(ctx context.Context, id string, version int) (entity.Bid, error)
	SendFeedback(ctx context.Context, req models.SendFeedback) error
	FindReviews(ctx context.Context, req models.FindReview) ([]entity.Review, error)
	SubmitApproveDecision(ctx context.Context, bidID, userID string) error
	GetBidApproveAmount(ctx context.Context, bidID string) (int, error)
	FindBidsByOrganization(ctx context.Context, organizationID string) ([]entity.Bid, error)
}
