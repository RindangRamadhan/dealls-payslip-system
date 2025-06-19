package v1

import (
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	logger           logger.Interface
	validator        *validator.Validate
	user             usecase.User
	payroll          usecase.Payroll
	payslip          usecase.Payslip
	overtime         usecase.Overtime
	attendance       usecase.Attendance
	reimbursement    usecase.Reimbursement
	attendancePeriod usecase.AttendancePeriod
}
