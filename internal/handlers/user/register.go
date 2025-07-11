package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"otus/go-server-project/internal/models"
)

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		u  models.User
		id string
	)
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id, err = h.service.RegisterUser(u); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"user_id":"%s"}`, id)))
}
