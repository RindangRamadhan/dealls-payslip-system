package payslip

import (
	"context"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/response"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/google/uuid"
)

type PayslipUsecase struct {
	logger            logger.Logger
	repo              repo.PayslipRepo
	payrollRepo       repo.PayrollRepo
	overtimeRepo      repo.OvertimeRepo
	attendanceRepo    repo.AttendanceRepo
	reimbursementRepo repo.ReimbursementRepo
}

// NewPayslipUsecase -.
func NewPayslipUsecase(
	l logger.Logger,
	r repo.PayslipRepo,
	pr repo.PayrollRepo,
	or repo.OvertimeRepo,
	ar repo.AttendanceRepo,
	rp repo.ReimbursementRepo,
) *PayslipUsecase {
	return &PayslipUsecase{
		logger:            l,
		repo:              r,
		payrollRepo:       pr,
		overtimeRepo:      or,
		attendanceRepo:    ar,
		reimbursementRepo: rp,
	}
}

func (uc *PayslipUsecase) GetPayslip(ctx context.Context, userID, periodID uuid.UUID) (*response.PayslipResponse, error) {
	// Get payroll for the period
	payroll, err := uc.payrollRepo.GetPayrollByPeriod(ctx, periodID)
	if err != nil {
		uc.logger.Error(err, "PayslipUsecase - GetPayslip - uc.payrollRepo.GetPayrollByPeriod")
		return nil, err
	}

	// Get user's payslip
	payslip, err := uc.repo.GetPayslipByUserAndPayroll(ctx, userID, payroll.ID)
	if err != nil {
		uc.logger.Error(err, "PayslipUsecase - GetPayslip - uc.repo.GetPayslipByUserAndPayroll")
		return nil, err
	}

	// Get supporting records
	attendanceRecords, err := uc.attendanceRepo.GetAttendanceByUserAndPeriod(ctx, userID, periodID)
	if err != nil {
		attendanceRecords = []entity.Attendance{}
	}

	overtimeRecords, err := uc.overtimeRepo.GetOvertimeByUserAndPeriod(ctx, userID, periodID)
	if err != nil {
		overtimeRecords = []entity.Overtime{}
	}

	reimbursementRecords, err := uc.reimbursementRepo.GetReimbursementsByUserAndPeriod(ctx, userID, periodID)
	if err != nil {
		reimbursementRecords = []entity.Reimbursement{}
	}

	return &response.PayslipResponse{
		Payslip:        *payslip,
		Attendances:    attendanceRecords,
		Overtimes:      overtimeRecords,
		Reimbursements: reimbursementRecords,
	}, nil
}
