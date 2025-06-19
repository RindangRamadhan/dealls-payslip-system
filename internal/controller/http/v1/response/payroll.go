package response

import "github.com/RindangRamadhan/dealls-payslip-system/internal/entity"

type PayrollSummaryResponse struct {
	Payroll      entity.Payroll `json:"payroll"`
	EmployeePays []EmployeePay  `json:"employee_pays"`
	TotalAmount  float64        `json:"total_amount"`
}

type (
	EmployeePay struct {
		UserID   string  `json:"user_id"`
		Username string  `json:"username"`
		TotalPay float64 `json:"total_pay"`
	}
)

func NewPayrollSummaryResponse(payroll *entity.Payroll, employeePays []EmployeePay) *PayrollSummaryResponse {
	return &PayrollSummaryResponse{
		Payroll:      *payroll,
		EmployeePays: employeePays,
		TotalAmount:  payroll.TotalAmount,
	}
}
