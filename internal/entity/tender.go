package entity

import (
	"time"

	"github.com/invopop/validation"
)

// TenderStatus is enum that represents all possible tender statuses.
type TenderStatus string

func (t TenderStatus) ValidationRule() validation.Rule {
	return validation.In(
		TenderCreated,
		TenderClosed,
		TenderPublished,
	)
}

func (t TenderStatus) String() string {
	return string(t)
}

const (
	TenderCreated   TenderStatus = "Created"
	TenderPublished TenderStatus = "Published"
	TenderClosed    TenderStatus = "Closed"
)

// ServiceType is enum that represents all possible service types.
type ServiceType string

func (t ServiceType) ValidationRule() validation.Rule {
	return validation.In(
		ServiceDelivery,
		ServiceConstruction,
		ServiceManufacture,
	)
}

func (t ServiceType) String() string {
	return string(t)
}

const (
	ServiceConstruction ServiceType = "Construction"
	ServiceDelivery     ServiceType = "Delivery"
	ServiceManufacture  ServiceType = "Manufacture"
)

// Tender is the entity that represents tender.
type Tender struct {
	ID              string       `json:"id" db:"id"`
	Name            string       `json:"name" db:"name"`
	Description     string       `json:"description" db:"description"`
	ServiceType     ServiceType  `json:"serviceType" db:"service_type"`
	Status          TenderStatus `json:"status" db:"status"`
	OrganizationID  string       `json:"organizationId" db:"organization_id"`
	CreatorUsername string       `json:"creatorUsername" db:"creator_username"`
	CreatedAt       time.Time    `json:"createdAt" db:"created_at"`
	Version         int          `json:"version" db:"version"`
}
