package middlewares

import (
	"net/http"
	"otus/go-server-project/internal/service"
)

// AuthMiddleware checks for Bearer token in Authorization header.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authenticator := service.Authenticator{}
		token, err := authenticator.ValidateToken(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Add("X-Authenticated-User", token)
		next.ServeHTTP(w, r)
	})
}
