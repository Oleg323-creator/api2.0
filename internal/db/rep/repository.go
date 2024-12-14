package rep

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

//USING THIS FILE FOR GETTING DATA FROM DB FOR GET/POST REQUESTS

// PARAMS FROM REQUESTS
type FilterParams struct {
	FromCurrency string  `form:"from_currency" json:"from_currency"`
	ToCurrency   string  `form:"to_currency" json:"to_currency"`
	Provider     string  `form:"provider" json:"provider"`
	Page         int     `form:"page" json:"page"`
	Limit        int     `form:"limit" json:"limit"`
	ID           int     `form:"id" json:"id"`
	Rate         float64 `form:"rate" json:"rate"`
	Order        string  `form:"order" json:"order"`
	OrderDir     string  `form:"order_dir" json:"order_dir"`
	Amount       float64 `form:"amount" json:"amount"`
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

// SAIVING THIS INFO INTO DB
type User struct {
	ID       int
	Email    string
	Password string
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
func (r *Repository) SignUpUserInDB(userData User) error {

	queryBuilder := squirrel.Insert("users").
		Columns("email", "password").
		Values(userData.Email, userData.Password)

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("Failed to build insert query: %v ", err)
	}

	// INSEARTING
	_, err = r.DB.Exec(query, args...)
	if err != nil {
		// CHECKING IF UNIQUE
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return errors.New("user already exists")
		}
		return err
	}
	return nil
}

// ERRORS FOR EXISTS ACCOUNT
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrEmailNotFound = errors.New("email not found")
var ErrInvalidPassword = errors.New("invalid password")

func (r *Repository) SignInUserInDB(email string) (string, string, error) {

	var storedEmail, storedPassword string

	query, args, err := squirrel.Select("email", "password").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", "", err
	}

	err = r.DB.QueryRow(query, args...).Scan(&storedEmail, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", ErrEmailNotFound
		}
		return "", "", err
	}

	return storedEmail, storedPassword, nil
}

func (r *Repository) SaveDataToDB(rateFrom string, rateTo string, connectorType string, data map[string]interface{}) error {

	// CHECKING KEYS WE ARE GETTING FROM PROVIDERS
	rate, ok := data[""]
	if !ok {
		rate, ok = data["USDT"]
		if !ok {
			rate, ok = data["BTC"]
			if !ok {
				rate, ok = data["ETH"]
				if !ok {
					rate, ok = data["BNB"]
					if !ok {
						return fmt.Errorf("invalid rate data format: got %T", data["rate"])
					}
				}
			}
		}
	}

	queryBuilder := squirrel.Insert("rates").
		Columns("from_currency", "to_currency", "rate", "provider", "created_at", "updated_at").
		Values(rateFrom, rateTo, rate, connectorType, time.Now(), time.Now()).
		Suffix("ON CONFLICT (from_currency, to_currency, provider) DO UPDATE SET rate = EXCLUDED.rate, updated_at = EXCLUDED.updated_at")

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Data saved to DB:", data)

	return nil
}
