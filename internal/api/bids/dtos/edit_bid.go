package dtos

import "github.com/invopop/validation"

type EditBidBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EditBidRequest struct {
	BidID    string `json:"bidId"`
	Username string `json:"username"`
	EditBidBody
}

func (r EditBidRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BidID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Name, validation.Length(0, 100)),
		validation.Field(&r.Description, validation.Length(0, 500)))
}
