package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendancePeriod struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by" db:"updated_by"`
}
