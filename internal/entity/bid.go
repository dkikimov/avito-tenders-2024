package entity

import (
	"time"

	"github.com/invopop/validation"
)

type BidStatus string

func (s BidStatus) ValidationRule() validation.Rule {
	return validation.In(
		BidCreated,
		BidPublished,
		BidCanceled,
		BidApproved,
		BidRejected)
}

const (
	BidCreated   BidStatus = "Created"
	BidPublished BidStatus = "Published"
	BidCanceled  BidStatus = "Canceled"
	BidApproved  BidStatus = "Approved"
	BidRejected  BidStatus = "Rejected"
)

type AuthorType string

func (t AuthorType) ValidationRule() validation.Rule {
	return validation.In(
		AuthorOrganization,
		AuthorUser)
}

const (
	AuthorOrganization AuthorType = "Organization"
	AuthorUser         AuthorType = "User"
)

type Decision string

func (s Decision) ValidationRule() validation.Rule {
	return validation.In(
		DecisionApproved,
		DecisionRejected,
	)
}

const (
	DecisionApproved = "Approved"
	DecisionRejected = "Rejected"
)

type Bid struct {
	Id          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Status      BidStatus  `json:"status" db:"status"`
	TenderId    string     `json:"tenderId" db:"tender_id"`
	AuthorType  AuthorType `json:"authorType" db:"author_type"`
	AuthorId    string     `json:"authorId" db:"author_id"`
	Version     int        `json:"version" db:"version"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
}
