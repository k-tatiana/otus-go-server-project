package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"otus/go-server-project/internal"
	"otus/go-server-project/internal/handlers/dialog"
	"otus/go-server-project/internal/handlers/post"
	"otus/go-server-project/internal/handlers/user"
	"otus/go-server-project/internal/handlers/websocket"
	"otus/go-server-project/internal/middlewares"
	"otus/go-server-project/internal/service"
	"otus/go-server-project/internal/transport/storage"
	"otus/go-server-project/internal/transport/storage/cache"
	"otus/go-server-project/internal/transport/storage/repository"
)

const (
	cacheCapacity = 10000
	useCache      = false
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

	defaultRouter := mux.NewRouter()

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

	repo, err := repository.NewRepo(ctx, masterDbCfg, slaveDbCfgs, env.DB.UseReplicas)
	if err != nil {
		log.Fatalf("Could not create repository: %v", err)
	}
	hasher := service.NewSimpleHasher(env.Secret)
	authenticator := service.NewAuthenticator(repo)

	userService := service.NewUserService(repo, hasher)
	userHandler := user.NewUserHandler(userService, logger, authenticator)

	cache := cache.NewLRUCache(cacheCapacity)
	cache.Prepare()

	strg := storage.NewStorage(repo, cache, env.Cache.UseCache)

	postsService := service.NewPostsService(strg)
	postsHandler := post.NewPostsHandler(postsService, logger, authenticator)

	dialogService := service.NewDialogService(repo)
	dialogHandler := dialog.NewDialogHandler(dialogService, logger, authenticator)

	r := defaultRouter.PathPrefix("/api").Subrouter()
	a := r.NewRoute().Subrouter() // for requests only for logged-in

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
	a.HandleFunc("/post/create", postsHandler.CreatePost).Methods("POST")
	// r.HandleFunc("/post/update", post.UpdatePost).Methods("PUT")
	// r.HandleFunc("/post/delete/{id}", post.DeletePost).Methods("PUT")
	// r.HandleFunc("/post/get/{id}", post.GetPost).Methods("GET")
	a.HandleFunc("/post/feed", postsHandler.FeedPost).Methods("GET")

	// // Dialog
	a.HandleFunc("/dialog/{user_id}/send", dialogHandler.SendDialog).Methods("POST")
	a.HandleFunc("/dialog/{user_id}/list", dialogHandler.ListDialog).Methods("GET")

	websocketHandler := websocket.NewWebSocketHandler(authenticator)
	go websocketHandler.Hub.Run()
	defaultRouter.HandleFunc("/ws", websocketHandler.HandleWebSocket)

	// middlewares
	a.Use(middlewares.AuthMiddleware)
	r.Use(middlewares.Logger)
	r.Use(middlewares.Responses)
	r.Use(middlewares.RequestIDMiddleware)

	s.srv.Handler = defaultRouter

	return s.srv.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.srv.Shutdown(ctx)
}
