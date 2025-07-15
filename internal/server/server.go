package server

import (
	"context"
	"log"
	"net/http"
	"otus/go-server-project/internal"
	"otus/go-server-project/internal/handlers/user"
	"otus/go-server-project/internal/middlewares"
	"otus/go-server-project/internal/repository"
	"otus/go-server-project/internal/service"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type HttpServer struct {
	srv *http.Server
}

func NewServer(addr string) *HttpServer {
	return &HttpServer{
		srv: &http.Server{
			Addr: addr,
		},
	}
}

func (s *HttpServer) Start() error {
	log.Printf("Starting server on %s", s.srv.Addr)

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	env, err := internal.EnvParse()
	if err != nil {
		log.Fatalf("Could not parse environment variables: %v", err)
	}

	r := mux.NewRouter()

	pool, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     env.DB.Host,
				Port:     uint16(env.DB.Port),
				User:     env.DB.User,
				Password: env.DB.Password,
				Database: env.DB.Database,
			},
			MaxConnections: 2,
		},
	)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	repo := repository.NewRepo(pool)
	hasher := service.NewSimpleHasher(env.Secret)

	userService := service.NewUserService(repo, hasher)
	userHandler := user.NewUserHandler(userService)

	a := r.NewRoute().Subrouter()

	// User
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/user/register", userHandler.RegisterUser).Methods("POST")
	a.HandleFunc("/user/get/{id}", userHandler.GetUser).Methods("GET")
	// r.HandleFunc("/user/search", user.SearchUser).Methods("GET")

	// // Friend
	// r.HandleFunc("/friend/set/{user_id}", friend.SetFriend).Methods("PUT")
	// r.HandleFunc("/friend/delete/{user_id}", friend.DeleteFriend).Methods("PUT")

	// // Post
	// r.HandleFunc("/post/create", post.CreatePost).Methods("POST")
	// r.HandleFunc("/post/update", post.UpdatePost).Methods("PUT")
	// r.HandleFunc("/post/delete/{id}", post.DeletePost).Methods("PUT")
	// r.HandleFunc("/post/get/{id}", post.GetPost).Methods("GET")
	// r.HandleFunc("/post/feed", post.FeedPost).Methods("GET")

	// // Dialog
	// r.HandleFunc("/dialog/{user_id}/send", dialog.SendDialog).Methods("POST")
	// r.HandleFunc("/dialog/{user_id}/list", dialog.ListDialog).Methods("GET")

	// middlewares
	a.Use(middlewares.AuthMiddleware)
	r.Use(middlewares.Logger)
	r.Use(middlewares.Responses)

	s.srv.Handler = r

	return s.srv.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.srv.Shutdown(ctx)
}
