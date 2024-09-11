package entity

import "time"

type OrganizationType string

const (
	IE  OrganizationType = "IE"
	LLC OrganizationType = "LLC"
	JSC OrganizationType = "JSC"
)

type Organization struct {
	Id          string           `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Type        OrganizationType `json:"type" db:"type"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

type OrganizationResponsible struct {
	Id             string
	OrganizationId int
	UserId         int
}
