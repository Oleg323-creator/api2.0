package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"net/http"
	"strconv"
)

// TO CONNECT IT WITH db.repository
type Handler struct {
	repository *db.Repository
}

func NewHandler(repository *db.Repository) *Handler {
	return &Handler{repository: repository}
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
