package handlers

import (
	_ "github.com/Masterminds/squirrel"
	"github.com/Oleg323-creator/api2.0/internal/db/rep"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GETTING THIS STRUCT FROM POST REQUEST
type CalculationRequest struct {
	FromCurrency string  `form:"from_currency" json:"from_currency"`
	ToCurrency   string  `form:"to_currency" json:"to_currency"`
	Provider     string  `form:"provider" json:"provider"`
	Amount       float64 `form:"amount" json:"amount"`
}

func (h *Handler) PostEndpoint(c *gin.Context) {
	var req CalculationRequest

	log.Println("Received POST request")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Decoded request: %+v", req)

	params := rep.FilterParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Provider:     req.Provider,
		Amount:       req.Amount,
	}

	// GO TO DB
	data, err := h.repository.GetRatesToCount(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data successfully fetched",
		"data":    data,
	})
}
