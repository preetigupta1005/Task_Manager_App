package models

import "time"

type Todo struct {
	Id          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	IsCompleted bool       `json:"isCompleted" db:"is_completed"`
	UserID      string     `json:"userId" db:"user_id"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt  *time.Time `json:"archivedAt" db:"archived_at"`
}

type TodoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	UserID      string `json:"userId" db:"user_id"`
}
