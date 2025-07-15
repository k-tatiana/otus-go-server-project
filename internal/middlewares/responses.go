package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const RetryAfter = 5 // seconds

type errorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Responses(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		if rw.status == http.StatusInternalServerError {
			resp := errorResponse{
				Message:   "internal server error",
				RequestID: requestID,
				Code:      500,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}
		if rw.status == http.StatusServiceUnavailable {
			resp := errorResponse{
				Message:   "service unavailable",
				RequestID: requestID,
				Code:      503,
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", fmt.Sprintf("%d", RetryAfter))
			json.NewEncoder(w).Encode(resp)
		}
		log.Printf("Response status: %d, Request ID: %s", rw.status, requestID)
	})
}
