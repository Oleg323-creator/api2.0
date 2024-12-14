package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

// FOR CHECKING ACCESS
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		if !isValidToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
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

/*

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

*/
