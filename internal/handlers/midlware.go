package handlers

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// FOR LOGS
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

// FOR CHECKING ACCESS
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		if !isValidToken(token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidToken(tokenString string) bool {
	// DEL "Bearer " FROM TOKEN
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}
