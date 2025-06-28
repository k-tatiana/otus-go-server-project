package handlers

import "net/http"

func SendDialog(w http.ResponseWriter, r *http.Request) {
	// TODO: implement send message logic
	w.WriteHeader(http.StatusOK)
}
