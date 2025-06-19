package request

type AttendancePeriodRequest struct {
	StartDate string `json:"start_date" validate:"required" example:"2025-01-01"`
	EndDate   string `json:"end_date" validate:"required" example:"2025-01-31"`
}
