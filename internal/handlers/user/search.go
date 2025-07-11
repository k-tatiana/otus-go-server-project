package user

import "net/http"

func SearchUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement search logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}
