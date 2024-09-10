package entity

import "time"

type TenderStatus string

func (t TenderStatus) String() string {
	return string(t)
}

const (
	TenderCreated   TenderStatus = "CREATED"
	TenderPublished TenderStatus = "PUBLISHED"
	TenderClosed    TenderStatus = "CLOSED"
)

// Tender is the entity that represents tender.
type Tender struct {
	Id              int          `json:"id" db:"id"`
	Name            string       `json:"name" db:"name"`
	Description     string       `json:"description" db:"description"`
	ServiceType     string       `json:"serviceType" db:"service_type"`
	Status          TenderStatus `json:"status" db:"status"`
	OrganizationId  int          `json:"organizationId" db:"organization_id"`
	CreatorUsername string       `json:"creatorUsername" db:"creator_username"`
	CreatedAt       time.Time    `json:"createdAt" db:"created_at"`
	Version         int          `json:"version" db:"version"`
}
