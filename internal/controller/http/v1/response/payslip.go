package response

import "github.com/RindangRamadhan/dealls-payslip-system/internal/entity"

type PayslipResponse struct {
	Payslip        entity.Payslip         `json:"payslip"`
	Attendances    []entity.Attendance    `json:"attendances"`
	Overtimes      []entity.Overtime      `json:"overtimes"`
	Reimbursements []entity.Reimbursement `json:"reimbursements"`
}
