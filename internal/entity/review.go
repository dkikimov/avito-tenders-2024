package entity

import "time"

type Review struct {
	Id          string    `db:"id"`
	Description string    `db:"description"`
	BidId       string    `db:"bid_id"`
	CreatedAt   time.Time `db:"created_at"`
}
