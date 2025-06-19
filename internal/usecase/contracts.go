// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/response"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=contracts.go -destination=mocks/mock_usecase.go -package=mock
type (
	User interface {
		GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	}

	Attendance interface {
		SubmitAttendance(ctx context.Context, req request.AttendanceRequest) error
	}

	AttendancePeriod interface {
		CreateAttendancePeriod(ctx context.Context, req *entity.AttendancePeriod) error
	}

	Overtime interface {
		SubmitOvertime(ctx context.Context, req request.OvertimeRequest) error
	}

	Payroll interface {
		GetPayrollSummary(ctx context.Context, periodID uuid.UUID) (*response.PayrollSummaryResponse, error)
		ProcessPayroll(ctx context.Context, periodID, userID uuid.UUID) (*entity.Payroll, error)
	}

	Payslip interface {
		GetPayslip(ctx context.Context, userID, periodID uuid.UUID) (*response.PayslipResponse, error)
	}

	Reimbursement interface {
		SubmitReimbursement(ctx context.Context, req request.ReimbursementRequest) error
	}
)
