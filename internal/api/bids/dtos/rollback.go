package dtos

import "github.com/invopop/validation"

type RollbackRequest struct {
	BidId    string `json:"bidId"`
	Version  int    `json:"version"`
	Username string `json:"username"`
}

func (r RollbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidId, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Version, validation.Required, validation.Min(1)),
		validation.Field(&r.Username, validation.Required))
}
