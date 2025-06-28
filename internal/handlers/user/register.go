package handlers

import "net/http"

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement registration logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"user_id":"example-user-id"}`))
}
