package handlers

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	_ "github.com/Masterminds/squirrel"
	"github.com/Oleg323-creator/api2.0/internal/db"
	_ "log"
	"net/http"
	"strconv"
)

type Handler struct {
	repository *db.Repository
}

// NewHandler создает новый экземпляр UserHandler.
func NewHandler(repository *db.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetEndpoint(w http.ResponseWriter, r *http.Request) {

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 3
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // По умолчанию устанавливаем страницу 1.
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

/*
func getEndpoint(w http.ResponseWriter, r *http.Request) {
	query := sq.
		Select("from_currency", "to_currency","rates", "provider").
		From("rates").

	// Генерация SQL-запроса и аргументов
	sqlString, args, err := query.ToSql()
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение запроса
	rows, err := db.Query(sqlString, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var rates float64
		var from_currency, to_currency, provider string
		if err := rows.Scan(&from_currency, &to_currency, &rates, &provider); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("One %s costs %f %s,Pvider = %s", from_currency, rates, to_currency, provider)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}




*/ /*
func main() {
	// Подключение к базе данных (пример для PostgreSQL)
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем SQL-запрос с помощью Squirrel
	query := sq.
		Select("id", "name", "email").
		From("users").
		Where(sq.Eq{"status": "active"}).
		OrderBy("created_at DESC")

	// Генерация SQL-запроса и аргументов
	sqlString, args, err := query.ToSql()
	if err != nil {
		log.Fatal(err)
	}

	// Вывод сгенерированного SQL
	fmt.Println("SQL Query:", sqlString)
	fmt.Println("Arguments:", args)

	// Выполнение запроса
	rows, err := db.Query(sqlString, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", id, name, email)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// Получить список всех пользователей
func listUsers(w http.ResponseWriter, r *http.Request) {
	// Выполняем запрос на получение данных
	rows, err := dbPool.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Failed to scan user data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Если произошла ошибка при чтении строк
	if err := rows.Err(); err != nil {
		http.Error(w, "Error occurred while fetching data", http.StatusInternalServerError)
		return
	}

	// Устанавливаем тип контента и отправляем результат
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

/*
import (
	sq "github.com/Masterminds/squirrel"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/rates/", postEndpoint)    // Обработка POST-запроса
	http.HandleFunc("/rates/pair/", getEndpoint) // Обработка GET-запроса

	// Запуск HTTP-сервера
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

// Обработчик GET-запроса
func getEndpoint(w http.ResponseWriter, r *http.Request) {
	query := sq.
		Select("from_currency", "to_currency","rates", "provider").
		From("rates").

	// Генерация SQL-запроса и аргументов
	sqlString, args, err := query.ToSql()
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение запроса
	rows, err := db.Query(sqlString, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Чтение результатов
	for rows.Next() {
		var rates float64
		var from_currency, to_currency, provider string
		if err := rows.Scan(&from_currency, &to_currency, &rates, &provider); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("One %s costs %f %s,Pvider = %s", from_currency, rates, to_currency, provider)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// Обработчик POST-запроса
func postEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad request: unable to parse JSON", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Received: Name=%s, Message=%s\n", msg.Name, msg.Message)
}

*/

/*	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello, %s! Welcome to our service.\n", name)
}*/
