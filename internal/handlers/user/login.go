package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"otus/go-server-project/internal/models"
)

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: implement login logic
	var l models.Login

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(l.Login, l.Password)
	if err != nil {
		fmt.Printf("Login failed for user %s: %v\n", l.Login, err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"` + token + `"}`))
}
