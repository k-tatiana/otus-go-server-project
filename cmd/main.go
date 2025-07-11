package main

import (
	"log"
	"otus/go-server-project/internal/server"
)

func main() {
	srv := server.NewServer(":8080")

	if err := srv.Start(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
