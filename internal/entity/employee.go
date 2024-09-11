package entity

import "time"

type Employee struct {
	Id        string
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
