package handlers

import (
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type SignUpRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest

	log.Println("Received Sign Up request")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body" + err.Error()})
		return
	}

	log.Printf("Decoded Sign Up request: %+v", req)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password" + err.Error()})
		return
	}

	params := db.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = h.repository.SignUpUserInDB(params)
	if err != nil {
		log.Printf("Error Signing Up: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Signing Up:" + err.Error()})
		return
	}

	response := gin.H{
		"message": fmt.Sprintf("User %s created successfully", req.Email),
	}

	c.JSON(http.StatusCreated, response)
}
