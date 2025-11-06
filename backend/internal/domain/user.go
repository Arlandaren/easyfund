package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	Phone        string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
