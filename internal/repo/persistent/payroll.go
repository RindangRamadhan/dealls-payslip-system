package persistent

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type PayrollRepo struct {
	*postgres.Postgres
}

// PayrollRepo -.
func NewPayrollRepo(pg *postgres.Postgres) *PayrollRepo {
	return &PayrollRepo{pg}
}

func (r *PayrollRepo) CreatePayroll(ctx context.Context, payroll *entity.Payroll) error {
	sql, args, err := r.Builder.
		Insert("payrolls").
		Columns(
			"id",
			"period_id",
			"processed_at",
			"total_amount",
			"employee_count",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			payroll.ID,
			payroll.PeriodID,
			payroll.ProcessedAt,
			payroll.TotalAmount,
			payroll.EmployeeCount,
			payroll.CreatedAt,
			payroll.UpdatedAt,
			payroll.CreatedBy,
			payroll.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - CreatePayroll - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - CreatePayroll - Exec")
		return err
	}

	return nil
}

func (r *PayrollRepo) UpdateTotalAmountPayroll(ctx context.Context, id uuid.UUID, totalAmount float64) error {
	sql, args, err := r.Builder.
		Update("payrolls").
		Set("total_amount", totalAmount).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - UpdateTotalAmountPayroll - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - UpdateTotalAmountPayroll - Exec")
		return err
	}

	return nil
}

func (r *PayrollRepo) GetPayrollByPeriod(ctx context.Context, periodID uuid.UUID) (*entity.Payroll, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"period_id",
			"processed_at",
			"total_amount",
			"employee_count",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("payrolls").
		Where(squirrel.Eq{"period_id": periodID}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - GetPayrollByPeriod - ToSql")
		return nil, err
	}

	var payroll entity.Payroll
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&payroll.ID,
		&payroll.PeriodID,
		&payroll.ProcessedAt,
		&payroll.TotalAmount,
		&payroll.EmployeeCount,
		&payroll.CreatedAt,
		&payroll.UpdatedAt,
		&payroll.CreatedBy,
		&payroll.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - GetPayrollByPeriod - QueryRow")
		return nil, err
	}

	return &payroll, nil
}

func (r *PayrollRepo) CheckPayrollExists(ctx context.Context, periodID uuid.UUID) (bool, error) {
	sql, args, err := r.Builder.
		Select("COUNT(*)").
		From("payrolls").
		Where(squirrel.Eq{"period_id": periodID}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - CheckPayrollExists - ToSql")
		return false, err
	}

	var count int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		r.Logger.Error(err, "PayrollRepo - CheckPayrollExists - QueryRow")
		return false, err
	}

	return count > 0, nil
}
