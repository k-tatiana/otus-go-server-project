package post

import "net/http"

func FeedPost(w http.ResponseWriter, r *http.Request) {
	// TODO: implement feed logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}
