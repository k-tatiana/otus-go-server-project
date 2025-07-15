package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	err := h.service.ValidateToken(r.Header.Get("X-Authenticated-User"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := mux.Vars(r)["id"]
	user, err := h.service.Get(userID)
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
