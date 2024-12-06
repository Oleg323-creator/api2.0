package services

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
	http.HandleFunc("/rates/list", getEndpoint) // Обработка GET-запроса

	// Запуск HTTP-сервера
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
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

// Обработчик GET-запроса
func getEndpoint(w http.ResponseWriter, r *http.Request) {
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

/*	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	fmt.Fprintf(w, "Hello, %s! Welcome to our service.\n", name)
}*/
