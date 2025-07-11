package friend

import "net/http"

func SetFriend(w http.ResponseWriter, r *http.Request) {
	// TODO: implement set friend logic
	w.WriteHeader(http.StatusOK)
}
