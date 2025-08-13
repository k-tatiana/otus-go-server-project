package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"otus/go-server-project/internal/models"

	"github.com/gorilla/mux"
)

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.authenticator.Auth(ctx, r.Header.Get("X-Authenticated-User"))
	if err != nil && errors.Is(err, models.ErrUnauthorized) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Printf("Error validating token: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userID := mux.Vars(r)["id"]
	user, err := h.service.Get(ctx, userID)
	if err != nil {
		fmt.Printf("Unable to get UserID from database %w", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Convert user model to json %w", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
