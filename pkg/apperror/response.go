package apperror

import (
	"errors"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		http.Error(w, appErr.Message, appErr.Code)
		log.Printf("sending error response: %s", err)
		return
	}

	http.Error(w, "error occurred", http.StatusInternalServerError)
	log.Printf("Error occurred: %s", err)
}
