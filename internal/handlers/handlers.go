package handlers

import (
	"encoding/json"
	"fmt"
	_ "github.com/Masterminds/squirrel"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TO CONNECT IT WITH db.repository
type Handler struct {
	repository *db.Repository
}

func NewHandler(repository *db.Repository) *Handler {
	return &Handler{repository: repository}
}

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

func (h *Handler) GetEndpoint(w http.ResponseWriter, r *http.Request) {

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3 //DEFAULT VALUE
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 //DEFAULT VALUE
	}

	fromCurrency := r.URL.Query().Get("from_currency")
	toCurrency := r.URL.Query().Get("to_currency")
	provider := r.URL.Query().Get("provider")
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 1
	}
	rateStr := r.URL.Query().Get("rate")
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		rate = 0
	}
	order := r.URL.Query().Get("order")
	orderDir := r.URL.Query().Get("order_dir")

	params := db.FilterParams{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Provider:     provider,
		Page:         page,
		Limit:        limit,
		OrderDir:     orderDir,
		ID:           id,
		Rate:         rate,
		Order:        order,
	}

	data, err := h.repository.GetRatesFromDB(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

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

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
