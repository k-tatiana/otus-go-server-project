package dialog

import (
	"encoding/json"
	"errors"
	"net/http"
	"otus/go-server-project/internal/models"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *DialogHandler) ListDialog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := r.Header.Get("X-Authenticated-User")
	err := h.authenticator.Auth(ctx, currentUser)
	if err != nil && errors.Is(err, models.ErrUnauthorized) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		h.logger.Error("Error validating token: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userID := mux.Vars(r)["user_id"]
	dialogs, err := h.service.List(ctx, currentUser, userID)
	if err != nil {
		h.logger.Error("Error listing messages: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dialogs)
	if err != nil {
		h.logger.Error("Error encoding response: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("List dialog with user %s\n", zap.String("userId", userID))
	w.WriteHeader(http.StatusOK)
}
