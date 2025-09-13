package post

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type CreatePostBody struct {
	Text string `json:"text"`
}

func (p *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// authenticate the user
	currentUser := r.Header.Get("X-Authenticated-User")
	err := p.authenticator.Auth(ctx, currentUser)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostBody
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	postId, err := p.service.WritePost(ctx, currentUser, req.Text)
	if err != nil {
		p.logger.Error("Error creating post: %v\n", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(postId))
}
