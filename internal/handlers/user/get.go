package handlers

import "net/http"

func GetUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get user logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
