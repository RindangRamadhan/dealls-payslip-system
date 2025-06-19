package persistent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type PayslipRepo struct {
	*postgres.Postgres
}

// PayslipRepo -.
func NewPayslipRepo(pg *postgres.Postgres) *PayslipRepo {
	return &PayslipRepo{pg}
}

func (r *PayslipRepo) CreatePayslip(ctx context.Context, payslip *entity.Payslip) error {
	sql, args, err := r.Builder.
		Insert("payslips").
		Columns(
			"id",
			"user_id",
			"payroll_id",
			"period_id",
			"base_salary",
			"attendance_days",
			"working_days",
			"attendance_salary",
			"overtime_hours",
			"overtime_salary",
			"reimbursement_total",
			"total_pay",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			payslip.ID,
			payslip.UserID,
			payslip.PayrollID,
			payslip.PeriodID,
			payslip.BaseSalary,
			payslip.AttendanceDays,
			payslip.WorkingDays,
			payslip.AttendanceSalary,
			payslip.OvertimeHours,
			payslip.OvertimeSalary,
			payslip.ReimbursementTotal,
			payslip.TotalPay,
			payslip.CreatedAt,
			payslip.UpdatedAt,
			payslip.CreatedBy,
			payslip.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - CreatePayslip - ToSql")
		return err
	}

	fmt.Println("payslip.PayrollID", payslip.PayrollID)

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - CreatePayslip - Exec")
		return err
	}

	return nil
}

func (r *PayslipRepo) GetPayslipByUserAndPayroll(ctx context.Context, userID, payrollID uuid.UUID) (*entity.Payslip, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"user_id",
			"payroll_id",
			"period_id",
			"base_salary",
			"attendance_days",
			"working_days",
			"attendance_salary",
			"overtime_hours",
			"overtime_salary",
			"reimbursement_total",
			"total_pay",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("payslips").
		Where(squirrel.Eq{"user_id": userID, "payroll_id": payrollID}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - GetPayslipByUserAndPayroll - ToSql")
		return nil, err
	}

	var payslip entity.Payslip
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&payslip.ID,
		&payslip.UserID,
		&payslip.PayrollID,
		&payslip.PeriodID,
		&payslip.BaseSalary,
		&payslip.AttendanceDays,
		&payslip.WorkingDays,
		&payslip.AttendanceSalary,
		&payslip.OvertimeHours,
		&payslip.OvertimeSalary,
		&payslip.ReimbursementTotal,
		&payslip.TotalPay,
		&payslip.CreatedAt,
		&payslip.UpdatedAt,
		&payslip.CreatedBy,
		&payslip.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - GetPayslipByUserAndPayroll - QueryRow")
		return nil, err
	}

	return &payslip, nil
}

func (r *PayslipRepo) GetPayslipsByPayroll(ctx context.Context, payrollID uuid.UUID) ([]entity.Payslip, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"user_id",
			"payroll_id",
			"period_id",
			"base_salary",
			"attendance_days",
			"working_days",
			"attendance_salary",
			"overtime_hours",
			"overtime_salary",
			"reimbursement_total",
			"total_pay",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("payslips").
		Where(squirrel.Eq{"payroll_id": payrollID}).
		OrderBy("total_pay DESC").
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - GetPayslipsByPayroll - ToSql")
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "PayslipRepo - GetPayslipsByPayroll - Query")
		return nil, err
	}
	defer rows.Close()

	var payslips []entity.Payslip
	for rows.Next() {
		var p entity.Payslip

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.PayrollID,
			&p.PeriodID,
			&p.BaseSalary,
			&p.AttendanceDays,
			&p.WorkingDays,
			&p.AttendanceSalary,
			&p.OvertimeHours,
			&p.OvertimeSalary,
			&p.ReimbursementTotal,
			&p.TotalPay,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
			&p.UpdatedBy,
		)
		if err != nil {
			r.Logger.Error(err, "PayslipRepo - GetPayslipsByPayroll - Scan")
			return nil, err
		}

		payslips = append(payslips, p)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error(err, "PayslipRepo - GetPayslipsByPayroll - rows.Err")
		return nil, err
	}

	return payslips, nil
}
