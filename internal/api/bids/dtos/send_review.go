package dtos

import (
	"github.com/invopop/validation"
)

type SendReviewRequest struct {
	BidId       string `json:"bid"`
	BidFeedback string `json:"bidFeedback"`
	Username    string `json:"username"`
}

func (r SendReviewRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidId, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.BidFeedback, validation.Required, validation.Length(1, 1000)),
		validation.Field(&r.Username, validation.Required))
}
