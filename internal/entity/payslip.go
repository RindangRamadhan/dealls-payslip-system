package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payslip struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	PayrollID          uuid.UUID `json:"payroll_id" db:"payroll_id"`
	PeriodID           uuid.UUID `json:"period_id" db:"period_id"`
	BaseSalary         float64   `json:"base_salary" db:"base_salary"`
	AttendanceDays     int       `json:"attendance_days" db:"attendance_days"`
	WorkingDays        int       `json:"working_days" db:"working_days"`
	AttendanceSalary   float64   `json:"attendance_salary" db:"attendance_salary"`
	OvertimeHours      float64   `json:"overtime_hours" db:"overtime_hours"`
	OvertimeSalary     float64   `json:"overtime_salary" db:"overtime_salary"`
	ReimbursementTotal float64   `json:"reimbursement_total" db:"reimbursement_total"`
	TotalPay           float64   `json:"total_pay" db:"total_pay"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy          uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy          uuid.UUID `json:"updated_by" db:"updated_by"`
}
