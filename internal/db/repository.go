package db

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

//USING THIS FILE FOR GETTING DATA FROM DB

type RateInfo struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rates        float64 `json:"rates"`
	Provider     string  `json:"provider"`
}

type FilterParams struct {
	FromCurrency string `json:"from_Currency"`
	ToCurrency   string `json:"to_currency"`
	Provider     string `json:"provider"`
	Page         int    `json:"page"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetRatesFromDB(params FilterParams) ([]RateInfo, error) {

	const pageSize = 3

	offset := (params.Page - 1) * pageSize

	queryBuilder := squirrel.Select("from_currency", "to_currency", "rate", "provider").
		From("rates").
		Limit(uint64(pageSize)).
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

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	var rates []RateInfo

	// READING RESULTS
	for rows.Next() {
		var rate RateInfo
		if err = rows.Scan(&rate.FromCurrency, &rate.ToCurrency, &rate.Rates, &rate.Provider); err != nil {
			log.Fatal(err)
		}

		rates = append(rates, rate)
		fmt.Printf("One %s costs %f %s,Pvider = %s \n", rate.FromCurrency, rate.Rates, rate.ToCurrency, rate.Provider)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rates, nil
}
