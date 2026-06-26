package model

import "time"

type Todo struct {
	Id          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	IsCompleted bool       `json:"isCompleted" db:"is_completed"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	ArchivedAt  *time.Time `json:"archivedAt" db:"archived_at"`
}

type TodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
