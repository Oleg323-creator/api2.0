package db

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
	"strings"
)

//USING THIS FILE FOR GETTING DATA FROM DB FOR GET/POST REQUESTS

// PARAMS FROM REQUESTS
type FilterParams struct {
	FromCurrency string  `json:"from_Currency"`
	ToCurrency   string  `json:"to_currency"`
	Provider     string  `json:"provider"`
	Page         int     `json:"page"`
	Limit        int     `json:"limit"`
	ID           int     `json:"id"`
	Rate         float64 `json:"rate"`
	Order        string  `json:"order"`
	OrderDir     string  `json:"order_dir"`
	Amount       float64 `json:"amount"`
}

// TO CONNECT IT WITH DB
type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// RESULT OF GET REQUEST
type RateInfo struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float64 `json:"rate"`
	Provider     string  `json:"provider"`
	ID           int     `json:"id"`
}

// GET REQUESTS GETTING DATA FROM HERE
func (r *Repository) GetRatesFromDB(params FilterParams) ([]RateInfo, error) {

	offset := (params.Page - 1) * params.Limit

	queryBuilder := squirrel.Select("from_currency", "to_currency", "rate", "provider", "id").
		From("rates").
		Limit(uint64(params.Limit)).
		Offset(uint64(offset))

	if params.FromCurrency != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"from_currency": "%" + params.FromCurrency + "%"})
	}

	if params.ToCurrency != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"to_currency": "%" + params.ToCurrency + "%"})
	}

	if params.Provider != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"provider": "%" + params.Provider + "%"})
	}

	if params.Order != "" {
		orderDirection := strings.ToUpper(params.OrderDir)

		if orderDirection != "ASC" && orderDirection != "DESC" {
			orderDirection = "ASC"
		}

		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("%s %s", params.Order, orderDirection))
	} else {
		queryBuilder = queryBuilder.OrderBy("rate DESC") // Сортировка по умолчанию по ID в порядке убывания
	}

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	var rates []RateInfo

	// GETTING RESULTS
	for rows.Next() {
		var rate RateInfo
		if err = rows.Scan(&rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Provider, &rate.ID); err != nil {
			log.Fatal(err)
		}

		rates = append(rates, rate)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rates, nil
}

// RESULT OF POST REQUEST
type CountRate struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	CountedRate  float64 `json:"counted_rate"`
	Rate         float64 `json:"rate"`
}

// POST REQUESTS GETTING DATA FROM HERE
func (r *Repository) GetRatesToCount(countRateParams FilterParams) ([]CountRate, error) {

	queryBuilder := squirrel.Select("from_currency", "to_currency", "provider", "rate").
		From("rates")

	if countRateParams.Provider != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"provider": "%" + countRateParams.Provider + "%"})
	}

	if countRateParams.FromCurrency != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"from_currency": countRateParams.FromCurrency})
	}

	if countRateParams.ToCurrency != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"to_currency": countRateParams.ToCurrency})
	}

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}
	fmt.Println("Generated SQL Query: ", query)

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	var countedRates []CountRate

	// GETTING RESULTS
	for rows.Next() {
		var countedRate CountRate

		if err = rows.Scan(&countedRate.FromCurrency, &countedRate.ToCurrency,
			&countedRate.Provider, &countedRate.Rate); err != nil {
			log.Fatal(err)
		}

		countedRate.Amount = countRateParams.Amount
		countedRate.CountedRate = countedRate.Rate * countedRate.Amount

		countedRates = append(countedRates, countedRate)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return countedRates, nil
}
