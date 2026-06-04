package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *App) getTodos(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, err := a.DB.Query(`
		SELECT id, text, done, created_at, updated_at
		FROM todos
		ORDER BY id ASC
	`)
	if err != nil {
		log.Printf("failed to fetch todos: %v", err)
		http.Error(w, "failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Printf("failed to scan todo: %v", err)
			http.Error(w, "failed to scan todos", http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (a *App) createTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var input struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if input.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO todos (text)
		VALUES ($1)
		RETURNING id, text, done, created_at, updated_at
	`

	var todo Todo

	err := a.DB.QueryRow(query, input.Text).
		Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("failed to create todo: %v", err)
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (a *App) updateTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	parts := strings.Split(path, "/")

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	if len(parts) == 2 && parts[1] == "toggle" {
		query := `
			UPDATE todos
			SET done = NOT done, updated_at = NOW()
			WHERE id = $1
			RETURNING id, text, done, created_at, updated_at
		`

		var todo Todo
		err := a.DB.QueryRow(query, id).
			Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)

		if err != nil {
			log.Printf("failed to toggle todo: %v", err)
			http.Error(w, "failed to toggle todo", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
		return
	}

	var input struct {
		Text string `json:"text"`
		Done bool   `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if input.Text == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE todos
		SET text = $1, done = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, text, done, created_at, updated_at
	`

	var todo Todo
	err = a.DB.QueryRow(query, input.Text, input.Done, id).
		Scan(&todo.ID, &todo.Text, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		log.Printf("failed to update todo: %v", err)
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (a *App) deleteTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	_, err = a.DB.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		log.Printf("failed to delete todo: %v", err)
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}