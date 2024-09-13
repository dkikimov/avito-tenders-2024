package apperror

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type reasonResponse struct {
	Reason string `json:"reason"`
}

func SendError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		if err := json.NewEncoder(w).Encode(reasonResponse{Reason: appErr.Message}); err != nil {
			slog.Error("failed to send reason response", "error", err)
			http.Error(w, appErr.Message, appErr.Code)
			return
		}

		slog.Info("sending error response", "error", err)

		return
	}

	http.Error(w, "error occurred", http.StatusInternalServerError)
}
