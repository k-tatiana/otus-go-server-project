package post

import "net/http"

func GetPost(w http.ResponseWriter, r *http.Request) {
	// TODO: implement get post logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
