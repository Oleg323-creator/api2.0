package handlers

import (
	"log"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		rr := &ResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(w, r)

		log.Printf("Completed request: %s %s with status %d in %v",
			r.Method, r.URL.Path, rr.statusCode, time.Since(start))
	})
}
