package dtos

import "github.com/invopop/validation"

type RollbackTenderRequest struct {
	Username string `json:"username"`
	Version  int    `json:"version"`
}

func (r RollbackTenderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Version, validation.Required))
}

type RollbackTender struct {
	Username string `json:"username"`
	Version  int    `json:"version"`
}
