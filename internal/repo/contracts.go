// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks/mock_repo.go -package=mock

type (
	// ---------------------------------- Presistent Repositories ----------------------------------
	UserRepo interface {
		GetAllEmployees(ctx context.Context) ([]entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	}

	AttendancePeriodRepo interface {
		CreateAttendancePeriod(ctx context.Context, period *entity.AttendancePeriod) error
		GetActivePeriod(ctx context.Context) (*entity.AttendancePeriod, error)
		GetPeriodByID(ctx context.Context, id uuid.UUID) (*entity.AttendancePeriod, error)
	}

	AttendanceRepo interface {
		CreateAttendance(ctx context.Context, attendance *entity.Attendance) error
		UpdateCheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time, updatedBy uuid.UUID) error
		GetAttendanceByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Attendance, error)
		GetAttendanceByDate(ctx context.Context, userID, periodID uuid.UUID, date time.Time) (*entity.Attendance, error)
	}

	OvertimeRepo interface {
		CreateOvertime(ctx context.Context, overtime *entity.Overtime) error
		GetOvertimeByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Overtime, error)
		GetTotalOvertimeByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (float64, error)
	}

	ReimbursementRepo interface {
		CreateReimbursement(ctx context.Context, reimbursement *entity.Reimbursement) error
		GetReimbursementsByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Reimbursement, error)
	}

	PayrollRepo interface {
		CreatePayroll(ctx context.Context, payroll *entity.Payroll) error
		UpdateTotalAmountPayroll(ctx context.Context, id uuid.UUID, totalAmount float64) error
		GetPayrollByPeriod(ctx context.Context, periodID uuid.UUID) (*entity.Payroll, error)
		CheckPayrollExists(ctx context.Context, periodID uuid.UUID) (bool, error)
	}

	PayslipRepo interface {
		CreatePayslip(ctx context.Context, payslip *entity.Payslip) error
		GetPayslipByUserAndPayroll(ctx context.Context, userID, payrollID uuid.UUID) (*entity.Payslip, error)
		GetPayslipsByPayroll(ctx context.Context, payrollID uuid.UUID) ([]entity.Payslip, error)
	}

	AuditLogRepo interface {
		CreateAuditLog(ctx context.Context, log *entity.AuditLog) error
	}
	// ---------------------------------- Presistent Repositories ----------------------------------

	// ------------------------------------ WebAPI Repositories ------------------------------------
	// ------------------------------------ WebAPI Repositories ------------------------------------
)
