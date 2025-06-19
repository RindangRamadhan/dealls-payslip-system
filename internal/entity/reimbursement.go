package entity

import (
	"time"

	"github.com/google/uuid"
)

type Reimbursement struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	PeriodID    uuid.UUID `json:"period_id" db:"period_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy   uuid.UUID `json:"updated_by" db:"updated_by"`
}
