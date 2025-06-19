package entity

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	PeriodID  uuid.UUID  `json:"period_id" db:"period_id"`
	Date      time.Time  `json:"date" db:"date"`
	CheckIn   *time.Time `json:"check_in" db:"check_in"`
	CheckOut  *time.Time `json:"check_out" db:"check_out"`
	IPAddress string     `json:"ip_address" db:"ip_address"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
	UpdatedBy uuid.UUID  `json:"updated_by" db:"updated_by"`
}
