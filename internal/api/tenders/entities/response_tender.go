package entities

import "time"

type ResponseTender struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	ServiceType    string    `db:"service_type"`
	Status         string    `db:"status"`
	OrganizationId int       `db:"organization_id"`
	Version        int       `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
}
