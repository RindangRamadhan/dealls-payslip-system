package payroll

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/response"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/google/uuid"
)

type PayrollUsecase struct {
	logger               logger.Logger
	repo                 repo.PayrollRepo
	userRepo             repo.UserRepo
	payslipRepo          repo.PayslipRepo
	overtimeRepo         repo.OvertimeRepo
	attendanceRepo       repo.AttendanceRepo
	reimbursementRepo    repo.ReimbursementRepo
	attendancePeriodRepo repo.AttendancePeriodRepo
}

// NewPayrollUsecase -.
func NewPayrollUsecase(
	l logger.Logger,
	r repo.PayrollRepo,
	ur repo.UserRepo,
	pr repo.PayslipRepo,
	or repo.OvertimeRepo,
	ar repo.AttendanceRepo,
	rr repo.ReimbursementRepo,
	apr repo.AttendancePeriodRepo,
) *PayrollUsecase {
	return &PayrollUsecase{
		logger:               l,
		repo:                 r,
		userRepo:             ur,
		payslipRepo:          pr,
		overtimeRepo:         or,
		attendanceRepo:       ar,
		reimbursementRepo:    rr,
		attendancePeriodRepo: apr,
	}
}

func (uc *PayrollUsecase) GetPayrollSummary(ctx context.Context, periodID uuid.UUID) (*response.PayrollSummaryResponse, error) {
	payroll, err := uc.repo.GetPayrollByPeriod(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payroll by period: %w", err)
	}

	payslips, err := uc.payslipRepo.GetPayslipsByPayroll(ctx, payroll.ID)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - GetPayrollSummary - uc.payslipRepo.GetPayslipsByPayroll")
		return nil, err
	}

	var employeePays []response.EmployeePay
	for _, payslip := range payslips {
		user, err := uc.userRepo.GetUserByID(ctx, payslip.UserID)
		if err != nil {
			continue
		}

		employeePays = append(employeePays, response.EmployeePay{
			UserID:   user.ID.String(),
			Username: user.Username,
			TotalPay: payslip.TotalPay,
		})
	}

	return response.NewPayrollSummaryResponse(payroll, employeePays), nil
}

func (uc *PayrollUsecase) ProcessPayroll(ctx context.Context, periodID, userID uuid.UUID) (*entity.Payroll, error) {
	// Check if payroll already exists for this period
	exists, err := uc.repo.CheckPayrollExists(ctx, periodID)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - ProcessPayroll - uc.repo.CheckPayrollExists")
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("Payroll for this period already processed")
	}

	// Get attendance period
	period, err := uc.attendancePeriodRepo.GetPeriodByID(ctx, periodID)
	if err != nil {
		return nil, err
	}

	// Get all employees
	employees, err := uc.userRepo.GetAllEmployees(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate working days in the period
	workingDays := utils.CountWorkingDays(period.StartDate, period.EndDate)

	// Create payroll record
	payroll := &entity.Payroll{
		ID:            uuid.New(),
		PeriodID:      periodID,
		ProcessedAt:   time.Now(),
		TotalAmount:   0,
		EmployeeCount: len(employees),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}

	err = uc.repo.CreatePayroll(ctx, payroll)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - ProcessPayroll - uc.repo.CreatePayroll")
		return nil, err
	}

	var totalAmount float64

	for _, employee := range employees {
		payslip, err := uc.calculatePayslip(ctx, employee, period, payroll.ID, workingDays, userID)
		if err != nil {
			uc.logger.Error(err, "PayrollUsecase - ProcessPayroll - uc.calculatePayslip")
			return nil, err
		}

		jp, _ := json.Marshal(payslip)
		fmt.Println("--- Payslip Data ---", employee.Username, string(jp))

		err = uc.payslipRepo.CreatePayslip(ctx, payslip)
		if err != nil {
			uc.logger.Error(err, "PayrollUsecase - ProcessPayroll - uc.payslipRepo.CreatePayslip")
			return nil, err
		}

		totalAmount += payslip.TotalPay
	}

	err = uc.repo.UpdateTotalAmountPayroll(ctx, payroll.ID, totalAmount)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - ProcessPayroll - uc.repo.UpdateTotalAmountPayroll")
		return nil, err
	}

	return nil, err
}

func (uc *PayrollUsecase) calculatePayslip(ctx context.Context, employee entity.User, period *entity.AttendancePeriod, payrollID uuid.UUID, workingDays int, userID uuid.UUID) (*entity.Payslip, error) {
	// Get attendance records
	attendances, err := uc.attendanceRepo.GetAttendanceByUserAndPeriod(ctx, employee.ID, period.ID)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - calculatePayslip - uc.attendanceRepo.GetAttendanceByUserAndPeriod")
		return nil, err
	}

	// Get overtime records
	overtimes, err := uc.overtimeRepo.GetOvertimeByUserAndPeriod(ctx, employee.ID, period.ID)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - calculatePayslip - uc.overtimeRepo.GetOvertimeByUserAndPeriod")
		return nil, err
	}

	// Get reimbursements
	reimbursements, err := uc.reimbursementRepo.GetReimbursementsByUserAndPeriod(ctx, employee.ID, period.ID)
	if err != nil {
		uc.logger.Error(err, "PayrollUsecase - calculatePayslip - uc.reimbursementRepo.GetReimbursementsByUserAndPeriod")
		return nil, err
	}

	// Calculate attendance salary
	attendanceDays := len(attendances)
	dailySalary := employee.Salary / float64(workingDays)
	attendanceSalary := dailySalary * float64(attendanceDays)

	// Calculate overtime salary
	var totalOvertimeHours float64
	for _, overtime := range overtimes {
		totalOvertimeHours += overtime.Hours
	}
	hourlyRate := dailySalary / 8                         // 8 hours per day
	overtimeSalary := hourlyRate * totalOvertimeHours * 2 // Double rate for overtime

	// Calculate total reimbursements
	var totalReimbursements float64
	for _, reimbursement := range reimbursements {
		totalReimbursements += reimbursement.Amount
	}

	// Calculate total pay
	totalPay := attendanceSalary + overtimeSalary + totalReimbursements

	payslip := &entity.Payslip{
		ID:                 uuid.New(),
		UserID:             employee.ID,
		PayrollID:          payrollID,
		PeriodID:           period.ID,
		BaseSalary:         employee.Salary,
		AttendanceDays:     attendanceDays,
		WorkingDays:        workingDays,
		AttendanceSalary:   attendanceSalary,
		OvertimeHours:      totalOvertimeHours,
		OvertimeSalary:     overtimeSalary,
		ReimbursementTotal: totalReimbursements,
		TotalPay:           totalPay,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          userID,
		UpdatedBy:          userID,
	}

	return payslip, nil
}
