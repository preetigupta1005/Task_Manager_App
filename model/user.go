package model

import "time"

type User struct {
	ID         string     `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Email      string     `db:"email" json:"email"`
	Password   string     `db:"password" json:"-"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	ArchivedAt *time.Time `db:"archived_at" json:"archivedAt"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=15"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=15"`
}

type LoginData struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
}
