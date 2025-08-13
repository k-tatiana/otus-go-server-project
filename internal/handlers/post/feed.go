package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (p *PostsHandler) FeedPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// authenticate the user
	err := p.authenticator.Auth(ctx, r.Header.Get("X-Authenticated-User"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Extract offset and limit from URL parameters
	offsetS := mux.Vars(r)["offset"]
	offset, err := strconv.Atoi(offsetS)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}
	limitS := mux.Vars(r)["limit"]
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}
	if limit <= 0 || offset < 0 {
		http.Error(w, "Invalid offset or limit", http.StatusBadRequest)
		return
	}
	p.logger.Info("Fetching feed", zap.Int("offset", offset), zap.Int("limit", limit))

	posts, err := p.service.Feed(ctx, offset, limit)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(posts) == 0 {
		http.Error(w, "No posts found", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Error converting posts to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
