package post

import "net/http"

func DeletePost(w http.ResponseWriter, r *http.Request) {
	// TODO: implement delete post logic
	w.WriteHeader(http.StatusOK)
}
