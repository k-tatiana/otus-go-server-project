package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Unable to get UserID from query params %w", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
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
