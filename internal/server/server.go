package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"otus/go-server-project/internal"
	"otus/go-server-project/internal/handlers/user"
	"otus/go-server-project/internal/middlewares"
	"otus/go-server-project/internal/repository"
	"otus/go-server-project/internal/service"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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

	ctx := context.Background()

	r := mux.NewRouter()

	logger := zap.NewNop()

	masterDbCfg := repository.Config(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			env.DBMaster.User,
			env.DBMaster.Password,
			env.DBMaster.Host,
			env.DBMaster.Port,
			env.DBMaster.Database,
		))
	slaveDbCfgs := []*pgxpool.Config{
		repository.Config(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s",
				env.DBReplica1.User,
				env.DBReplica1.Password,
				env.DBReplica1.Host,
				env.DBReplica1.Port,
				env.DBReplica1.Database,
			)),
		repository.Config(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s",
				env.DBReplica2.User,
				env.DBReplica2.Password,
				env.DBReplica2.Host,
				env.DBReplica2.Port,
				env.DBReplica2.Database,
			)),
	}

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	repo, err := repository.NewRepo(ctx, masterDbCfg, slaveDbCfgs)
	if err != nil {
		log.Fatalf("Could not create repository: %v", err)
	}
	hasher := service.NewSimpleHasher(env.Secret)

	userService := service.NewUserService(repo, hasher)
	userHandler := user.NewUserHandler(userService, logger)

	a := r.NewRoute().Subrouter()

	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	// User
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/user/register", userHandler.RegisterUser).Methods("POST")
	a.HandleFunc("/user/get/{id}", userHandler.GetUser).Methods("GET")
	a.HandleFunc("/user/search", userHandler.SearchUser).Methods("GET")

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
