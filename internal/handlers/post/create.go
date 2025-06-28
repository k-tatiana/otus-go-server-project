package handlers

import "net/http"

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// TODO: implement create post logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"id":"example-post-id"}`))
}
