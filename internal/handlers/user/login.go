package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"otus/go-server-project/internal/models"
	"otus/go-server-project/internal/service"
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
	if errors.Is(err, models.ErrNoSuchUser) {
		http.Error(w, "No such user", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if errors.Is(err, models.ErrInvalidCredentials) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		fmt.Printf("Login failed for user %s: %v\n", l.Login, err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	authenticator := service.Authenticator{}
	authToken := authenticator.GenerateToken(token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"` + authToken + `"}`))
}
