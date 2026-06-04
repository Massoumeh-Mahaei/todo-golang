package main

import "time"

type Todo struct {
	ID        int       `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	Done      bool      `json:"done" db:"done"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}