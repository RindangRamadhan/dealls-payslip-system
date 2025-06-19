package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payroll struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PeriodID      uuid.UUID `json:"period_id" db:"period_id"`
	ProcessedAt   time.Time `json:"processed_at" db:"processed_at"`
	TotalAmount   float64   `json:"total_amount" db:"total_amount"`
	EmployeeCount int       `json:"employee_count" db:"employee_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy     uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy     uuid.UUID `json:"updated_by" db:"updated_by"`
}
