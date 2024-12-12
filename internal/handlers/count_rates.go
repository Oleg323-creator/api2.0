package handlers

import (
	"encoding/json"
	"fmt"
	_ "github.com/Masterminds/squirrel"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"log"
	"net/http"
)

// GETTING THIS STRUCT FROM POST REQUEST
type CalculationRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
}

func (h *Handler) PostEndpoint(w http.ResponseWriter, r *http.Request) {
	var req CalculationRequest

	log.Println("Received POST request")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded request: %+v", req)

	params := db.FilterParams{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Provider:     req.Provider,
		Amount:       req.Amount,
	}

	// GO TO DB
	data, err := h.repository.GetRatesToCount(params)
	if err != nil {
		log.Printf("Error fetching rates: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching rates: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
