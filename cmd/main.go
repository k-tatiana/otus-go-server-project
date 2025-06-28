package main

import (
	"log"
	dialog "otus/go-server-project/internal/handlers/dialog"
	friend "otus/go-server-project/internal/handlers/friend"
	post "otus/go-server-project/internal/handlers/post"
	user "otus/go-server-project/internal/handlers/user"
	"otus/go-server-project/internal/server"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// User
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/user/register", user.RegisterUser).Methods("POST")
	r.HandleFunc("/user/get/{id}", user.GetUser).Methods("GET")
	r.HandleFunc("/user/search", user.SearchUser).Methods("GET")

	// Friend
	r.HandleFunc("/friend/set/{user_id}", friend.SetFriend).Methods("PUT")
	r.HandleFunc("/friend/delete/{user_id}", friend.DeleteFriend).Methods("PUT")

	// Post
	r.HandleFunc("/post/create", post.CreatePost).Methods("POST")
	r.HandleFunc("/post/update", post.UpdatePost).Methods("PUT")
	r.HandleFunc("/post/delete/{id}", post.DeletePost).Methods("PUT")
	r.HandleFunc("/post/get/{id}", post.GetPost).Methods("GET")
	r.HandleFunc("/post/feed", post.FeedPost).Methods("GET")

	// Dialog
	r.HandleFunc("/dialog/{user_id}/send", dialog.SendDialog).Methods("POST")
	r.HandleFunc("/dialog/{user_id}/list", dialog.ListDialog).Methods("GET")

	srv := server.NewServer(":8080", r)

	if err := srv.Start(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
