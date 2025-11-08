package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
