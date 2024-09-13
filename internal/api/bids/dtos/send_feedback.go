package dtos

import "github.com/invopop/validation"

type SendFeedbackRequest struct {
	BidID    string `json:"bidId"`
	Feedback string `json:"bidFeedback"`
	Username string `json:"username"`
}

func (r SendFeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Feedback, validation.Required, validation.Length(1, 1000)),
		validation.Field(&r.Username, validation.Required))
}
