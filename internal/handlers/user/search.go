package user

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (h *UserHandler) SearchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.authenticator.Auth(ctx, r.Header.Get("X-Authenticated-User"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	name := r.URL.Query().Get("firstName")
	surname := r.URL.Query().Get("secondName")

	users, err := h.service.SearchUser(ctx, name, surname)
	if err != nil {
		h.logger.Warn("SearchUser error", zap.String("name", name), zap.String("surname", surname))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if len(users) == 0 {
		h.logger.Warn("No users found", zap.String("name", name), zap.String("surname", surname))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		h.logger.Error("Convert user models to json", zap.Error(err))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
