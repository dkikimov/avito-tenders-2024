package entities

type EditTenderStatusRequest struct {
	Status   string `json:"status"`
	Username string `json:"username"`
}

type EditTenderRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ServiceType *string `json:"serviceType,omitempty"`
}
