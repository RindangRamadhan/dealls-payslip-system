package persistent

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type AttendancePeriodRepo struct {
	*postgres.Postgres
}

// AttendancePeriodRepo -.
func NewAttendancePeriodRepo(pg *postgres.Postgres) *AttendancePeriodRepo {
	return &AttendancePeriodRepo{pg}
}

func (r *AttendancePeriodRepo) CreateAttendancePeriod(ctx context.Context, period *entity.AttendancePeriod) error {
	sql, args, err := r.Builder.
		Insert("attendance_periods").
		Columns(
			"id",
			"start_date",
			"end_date",
			"is_active",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			period.ID,
			period.StartDate,
			period.EndDate,
			period.IsActive,
			period.CreatedAt,
			period.UpdatedAt,
			period.CreatedBy,
			period.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - CreateAttendancePeriod - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - CreateAttendancePeriod - Exec")
		return err
	}

	return nil
}

func (r *AttendancePeriodRepo) GetActivePeriod(ctx context.Context) (*entity.AttendancePeriod, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"start_date",
			"end_date",
			"is_active",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("attendance_periods").
		Where(squirrel.Eq{"is_active": true}).
		OrderBy("created_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - GetActivePeriod - ToSql")
		return nil, err
	}

	var period entity.AttendancePeriod
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&period.ID,
		&period.StartDate,
		&period.EndDate,
		&period.IsActive,
		&period.CreatedAt,
		&period.UpdatedAt,
		&period.CreatedBy,
		&period.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - GetActivePeriod - QueryRow")
		return nil, err
	}

	return &period, nil
}

func (r *AttendancePeriodRepo) GetPeriodByID(ctx context.Context, id uuid.UUID) (*entity.AttendancePeriod, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"start_date",
			"end_date",
			"is_active",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("attendance_periods").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - GetPeriodByID - ToSql")
		return nil, err
	}

	var period entity.AttendancePeriod
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&period.ID,
		&period.StartDate,
		&period.EndDate,
		&period.IsActive,
		&period.CreatedAt,
		&period.UpdatedAt,
		&period.CreatedBy,
		&period.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "AttendancePeriodRepo - GetPeriodByID - QueryRow")
		return nil, err
	}

	return &period, nil
}
