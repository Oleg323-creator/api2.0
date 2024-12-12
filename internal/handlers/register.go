package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	log.Println("Received Sign Up request")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded Sign Up request: %+v", req)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	params := db.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = h.repository.SignUpUserInDB(params)
	if err != nil {
		log.Printf("Error Signing Up: %v", err)
		http.Error(w, fmt.Sprintf("\"Error Signing Up: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := fmt.Sprintf(`{"message": "User %s created successfully"}`, req.Email)
	w.Write([]byte(response))
}
