package bids

import "net/http"

type HTTPHandlers interface {
	CreateBid(w http.ResponseWriter, r *http.Request)
	GetUserBids(w http.ResponseWriter, r *http.Request)
	RollbackTender(w http.ResponseWriter, r *http.Request)
}
