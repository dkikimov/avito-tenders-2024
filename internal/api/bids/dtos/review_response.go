package dtos

import (
	"avito-tenders/internal/entity"
	"avito-tenders/pkg/types"
)

type ReviewResponse struct {
	Id          string            `json:"id"`
	Description string            `json:"description"`
	CreatedAt   types.RFC3339Time `json:"createdAt"`
}

func NewReviewResponse(review entity.Review) ReviewResponse {
	return ReviewResponse{
		Id:          review.Id,
		Description: review.Description,
		CreatedAt:   types.RFCFromTime(review.CreatedAt),
	}
}

func NewReviewResponseList(reviews []entity.Review) []ReviewResponse {
	reviewResponseList := make([]ReviewResponse, 0, len(reviews))

	for i := range reviews {
		reviewResponseList = append(reviewResponseList, NewReviewResponse(reviews[i]))
	}

	return reviewResponseList
}
