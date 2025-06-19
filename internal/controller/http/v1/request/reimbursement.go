package request

import "github.com/google/uuid"

type ReimbursementRequest struct {
	Amount      float64 `json:"amount" validate:"required,min=0"`
	Description string  `json:"description" validate:"required"`

	UserID   uuid.UUID `json:"-"`
	ClientIP string    `json:"-"`
}
