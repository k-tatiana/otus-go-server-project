package dialogv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"otus/go-server-project/internal/models"
)

var req struct {
	Message string `json:"text"`
}

func (h *DialogAPIHandler) SendDialog(w http.ResponseWriter, r *http.Request) {
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
	h.logger.Info("Send dialog to user %s\n", zap.String("userId", userID))

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	message := req.Message
	if message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}
	h.logger.Info("Message: %s\n", zap.String("message", message))
	err = h.service.Send(ctx, authUserToken, userID, message, nil)
	if err != nil {
		h.logger.Error("Error sending message: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Message sent to user %s\n", userID)

	w.WriteHeader(http.StatusOK)
}
