package request

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceRequest struct {
	Date     string `json:"date" validate:"required" example:"2025-06-20"`
	CheckIn  string `json:"check_in,omitempty" example:"08:00:00"`
	CheckOut string `json:"check_out,omitempty" example:"17:00:00"`
	Type     string `json:"type" validate:"required,oneof=check_in check_out" example:"check_in"`

	UserID       uuid.UUID  `json:"-"`
	DateTime     time.Time  `json:"-"`
	CheckInTime  *time.Time `json:"-"`
	CheckOutTime *time.Time `json:"-"`
	ClientIP     string     `json:"-"`
}
