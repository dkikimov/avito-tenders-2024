package dtos

import "github.com/invopop/validation"

type SubmitDecisionRequest struct {
	BidId    string `json:"bidId"`
	Decision string `json:"decision"`
	Username string `json:"username"`
}

func (r SubmitDecisionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidId, validation.Required),
		validation.Field(&r.Decision, validation.Required))
}
