package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	Salary    float64   `json:"salary" db:"salary"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by" db:"updated_by"`
}
