package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
}
func main() {
	connStr := "host=postgres port=5432 user=postgres password=postgres dbname=tododb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	app := &App{DB: db}

http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getTodos(w, r)
	case http.MethodPost:
		app.createTodo(w, r)
	case http.MethodOptions:
		enableCORS(w)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
})

http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		app.updateTodo(w, r)
	case http.MethodDelete:
		app.deleteTodo(w, r)
	case http.MethodOptions:
		enableCORS(w)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}