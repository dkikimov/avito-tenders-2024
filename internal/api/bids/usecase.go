package bids

import (
	"context"

	"avito-tenders/internal/api/bids/dtos"
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/queryparams"
)

type Usecase interface {
	Create(ctx context.Context, req dtos.CreateBidRequest) (dtos.BidResponse, error)
	FindByUsername(ctx context.Context, username string, pagination queryparams.Pagination) ([]dtos.BidResponse, error)
	FindByTenderID(ctx context.Context, req dtos.FindByTenderIDRequest) ([]dtos.BidResponse, error)
	GetStatusByID(ctx context.Context, bidID string, username string) (entity.BidStatus, error)
	UpdateStatusByID(ctx context.Context, req dtos.UpdateStatusRequest) (dtos.BidResponse, error)
	Edit(ctx context.Context, req dtos.EditBidRequest) (dtos.BidResponse, error)
	SubmitDecision(ctx context.Context, req dtos.SubmitDecisionRequest) (dtos.BidResponse, error)
	SendFeedback(ctx context.Context, req dtos.SendFeedbackRequest) (dtos.BidResponse, error)
	Rollback(ctx context.Context, req dtos.RollbackRequest) (dtos.BidResponse, error)
	FindReviewsByTenderID(ctx context.Context, req dtos.FindReviewsRequest) ([]dtos.ReviewResponse, error)
}
