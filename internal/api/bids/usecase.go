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
	FindByTenderId(ctx context.Context, req dtos.FindByTenderIdRequest) ([]dtos.BidResponse, error)
	GetStatusById(ctx context.Context, bidId string, username string) (entity.BidStatus, error)
	UpdateStatusById(ctx context.Context, req dtos.UpdateStatusRequest) (dtos.BidResponse, error)
	SubmitDecision(ctx context.Context, req dtos.SubmitDecisionRequest) (dtos.BidResponse, error)
	SendFeedback(ctx context.Context, req dtos.SendFeedbackRequest) (dtos.BidResponse, error)
	Rollback(ctx context.Context, req dtos.RollbackRequest) (dtos.BidResponse, error)
	FindReviewsByTenderId(ctx, req dtos.FindReviewsRequest) ([]entity.Review, error)
}
