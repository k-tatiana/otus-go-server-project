package main

import (
	"log"

	_ "net/http/pprof"

	"otus/go-server-project/internal/server"
)

func main() {
	srv := server.NewServer(":8080")

	if err := srv.Start(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
