package entity

import "time"

type Review struct {
	ID          string    `db:"id"`
	Description string    `db:"description"`
	BidID       string    `db:"bid_id"`
	CreatedAt   time.Time `db:"created_at"`
}
