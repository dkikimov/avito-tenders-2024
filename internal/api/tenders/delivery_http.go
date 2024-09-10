package tenders

import "net/http"

type HTTPHandlers interface {
	CreateTender(w http.ResponseWriter, r *http.Request)
	GetUserTenders(w http.ResponseWriter, r *http.Request)
	EditTender(w http.ResponseWriter, r *http.Request)
	RollbackTender(w http.ResponseWriter, r *http.Request)
}
