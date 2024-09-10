package entity

type BidStatus string

const (
	BidCreated   BidStatus = "CREATED"
	BidPublished BidStatus = "PUBLISHED"
	BidCanceled  BidStatus = "CANCELED"
)

type Bid struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          BidStatus `json:"status"`
	TenderId        int       `json:"tenderId"`
	OrganizationId  int       `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
}
