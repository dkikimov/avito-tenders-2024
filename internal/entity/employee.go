package entity

import "time"

type Employee struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
