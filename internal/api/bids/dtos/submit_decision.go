package dtos

import (
	"github.com/invopop/validation"

	"avito-tenders/internal/entity"
)

type SubmitDecisionRequest struct {
	BidID    string             `json:"bidId"`
	Decision entity.BidDecision `json:"decision"`
	Username string             `json:"username"`
}

func (r SubmitDecisionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidID, validation.Required),
		validation.Field(&r.Decision, validation.Required, r.Decision.ValidationRule()),
		validation.Field(&r.Username, validation.Required))
}
