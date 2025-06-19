package request

type PayrollRequest struct {
	PeriodID string `json:"period_id" validate:"required" example:"1e6c4313-3148-4fc8-a6cd-b6cd4674e041"`
}
