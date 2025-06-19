package request

import (
	"time"

	"github.com/google/uuid"
)

type OvertimeRequest struct {
	Date  string  `json:"date" validate:"required" example:"2025-06-20"`
	Hours float64 `json:"hours" validate:"required,min=0.5,max=3" example:"1"`

	UserID   uuid.UUID `json:"-"`
	DateTime time.Time `json:"-"`
	ClientIP string    `json:"-"`
}
