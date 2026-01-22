package dialogv2

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"otus/go-server-project/internal/models"
)

func (h *DialogAPIHandler) GetDialogCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUserToken := r.Header.Get("X-Authenticated-User")
	err := h.authenticator.Auth(ctx, authUserToken)
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
	h.logger.Info("Get dialog count with user %s\n", zap.String("userId", userID))

	count, err := h.counter.GetCount(ctx, authUserToken, userID)
	if err != nil {
		h.logger.Error("Error getting dialog count: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
		h.logger.Error("Error encoding response: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("Dialog count with user %s retrieved successfully\n", zap.String("userId", userID))
	w.WriteHeader(http.StatusOK)
}
