package apperror

import (
	"errors"
	"log/slog"
	"net/http"
)

func SendError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		http.Error(w, appErr.Message, appErr.Code)
		slog.Info("sending error response", "error", err)
		return
	}

	http.Error(w, "error occurred", http.StatusInternalServerError)
}
