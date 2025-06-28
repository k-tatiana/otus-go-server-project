package handlers

import (
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: implement login logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"example-token"}`))
}
