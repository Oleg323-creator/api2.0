package handlers

import (
	"encoding/json"
	_ "errors"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWTtoken string `json:"token"`
}

var jwtSecretKey = []byte("sf9vd9s1vsfdvdsv8fdv56114869s5fvd1hntjmuhngrbvfretbhnymju")

// GenerateJWT
func GenerateJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// ADDING TO PAYLOAD
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	// SIGN TOKEN
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

	var req SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//GETTING USER DATA FROM DB
	storedUsername, storedPassword, err := h.repository.SignInUserInDB(req.Email)
	if err != nil {
		if err == db.ErrEmailNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//CHECKING PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(storedUsername)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := LoginResponse{
		JWTtoken: token,
	}
	json.NewEncoder(w).Encode(response)
}
