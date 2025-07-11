package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs all incoming HTTP requests.
func Logger(next http.Handler) http.Handler {
	fmt.Println("Logger middleware initialized")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logger middleware working")
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}
