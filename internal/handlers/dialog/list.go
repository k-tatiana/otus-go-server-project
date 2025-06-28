package handlers

import "net/http"

func ListDialog(w http.ResponseWriter, r *http.Request) {
	// TODO: implement list messages logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}
