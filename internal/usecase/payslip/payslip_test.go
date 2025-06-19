package payslip

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mock "github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/mocks"
)

func TestGetPayslip(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	periodID := uuid.New()
	payrollID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		repoPayslip := mock.NewMockPayslipRepo(ctrl)
		repoPayroll := mock.NewMockPayrollRepo(ctrl)
		repoOvertime := mock.NewMockOvertimeRepo(ctrl)
		repoAttendance := mock.NewMockAttendanceRepo(ctrl)
		repoReimbursement := mock.NewMockReimbursementRepo(ctrl)
		log := logger.New("debug")

		payroll := &entity.Payroll{ID: payrollID, PeriodID: periodID}
		payslipRecord := &entity.Payslip{
			ID:                 uuid.New(),
			PayrollID:          payrollID,
			UserID:             userID,
			PeriodID:           periodID,
			BaseSalary:         5000000,
			AttendanceDays:     20,
			WorkingDays:        22,
			AttendanceSalary:   4545454.54,
			OvertimeHours:      5,
			OvertimeSalary:     250000,
			ReimbursementTotal: 300000,
			TotalPay:           5095454.54,
			CreatedAt:          now,
			UpdatedAt:          now,
			CreatedBy:          userID,
			UpdatedBy:          userID,
		}

		repoPayroll.EXPECT().GetPayrollByPeriod(gomock.Any(), periodID).Return(payroll, nil)
		repoPayslip.EXPECT().GetPayslipByUserAndPayroll(gomock.Any(), userID, payrollID).Return(payslipRecord, nil)
		repoAttendance.EXPECT().GetAttendanceByUserAndPeriod(gomock.Any(), userID, periodID).Return([]entity.Attendance{}, nil)
		repoOvertime.EXPECT().GetOvertimeByUserAndPeriod(gomock.Any(), userID, periodID).Return([]entity.Overtime{}, nil)
		repoReimbursement.EXPECT().GetReimbursementsByUserAndPeriod(gomock.Any(), userID, periodID).Return([]entity.Reimbursement{}, nil)

		uc := NewPayslipUsecase(*log, repoPayslip, repoPayroll, repoOvertime, repoAttendance, repoReimbursement)

		res, err := uc.GetPayslip(context.Background(), userID, periodID)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, payslipRecord.ID, res.Payslip.ID)
	})

	t.Run("fail get payroll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoPayslip := mock.NewMockPayslipRepo(ctrl)
		repoPayroll := mock.NewMockPayrollRepo(ctrl)
		log := logger.New("debug")

		repoPayroll.EXPECT().GetPayrollByPeriod(gomock.Any(), periodID).Return(nil, errors.New("payroll not found"))

		uc := NewPayslipUsecase(*log, repoPayslip, repoPayroll, nil, nil, nil)
		res, err := uc.GetPayslip(context.Background(), userID, periodID)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("fail get payslip", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoPayslip := mock.NewMockPayslipRepo(ctrl)
		repoPayroll := mock.NewMockPayrollRepo(ctrl)
		log := logger.New("debug")

		payroll := &entity.Payroll{ID: payrollID, PeriodID: periodID}
		repoPayroll.EXPECT().GetPayrollByPeriod(gomock.Any(), periodID).Return(payroll, nil)
		repoPayslip.EXPECT().GetPayslipByUserAndPayroll(gomock.Any(), userID, payrollID).Return(nil, errors.New("payslip not found"))

		uc := NewPayslipUsecase(*log, repoPayslip, repoPayroll, nil, nil, nil)
		res, err := uc.GetPayslip(context.Background(), userID, periodID)
		require.Error(t, err)
		require.Nil(t, res)
	})
}
