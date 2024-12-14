package handlers

import (
	_ "errors"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type SignInRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type LoginResponse struct {
	JWTtoken string `form:"token" json:"token"`
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

func (h *Handler) SignIn(c *gin.Context) {

	var req SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body" + err.Error()})
		return
	}

	//GETTING USER DATA FROM DB
	storedUsername, storedPassword, err := h.repository.SignInUserInDB(req.Email)
	if err != nil {
		if err == db.ErrEmailNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found" + err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	//CHECKING PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password" + err.Error()})
		return
	}

	token, err := GenerateJWT(storedUsername)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	tokenJWT := LoginResponse{
		JWTtoken: token,
	}

	response := gin.H{
		"message": tokenJWT,
	}

	c.JSON(http.StatusCreated, response)
}

/*
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
*/
