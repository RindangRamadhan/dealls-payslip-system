package persistent

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type OvertimeRepo struct {
	*postgres.Postgres
}

// OvertimeRepo -.
func NewOvertimeRepo(pg *postgres.Postgres) *OvertimeRepo {
	return &OvertimeRepo{pg}
}

func (r *OvertimeRepo) CreateOvertime(ctx context.Context, overtime *entity.Overtime) error {
	sql, args, err := r.Builder.
		Insert("overtimes").
		Columns(
			"id",
			"user_id",
			"period_id",
			"date",
			"hours",
			"ip_address",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			overtime.ID,
			overtime.UserID,
			overtime.PeriodID,
			overtime.Date,
			overtime.Hours,
			overtime.IPAddress,
			overtime.CreatedAt,
			overtime.UpdatedAt,
			overtime.CreatedBy,
			overtime.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "OvertimeRepo - CreateOvertime - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "OvertimeRepo - CreateOvertime - Exec")
		return err
	}

	return nil
}

func (r *OvertimeRepo) GetOvertimeByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Overtime, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"user_id",
			"period_id",
			"date",
			"hours",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("overtimes").
		Where(squirrel.Eq{"user_id": userID, "period_id": periodID}).
		OrderBy("date").
		ToSql()
	if err != nil {
		r.Logger.Error(err, "OvertimeRepo - GetOvertimeByUserAndPeriod - ToSql")
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "OvertimeRepo - GetOvertimeByUserAndPeriod - Query")
		return nil, err
	}
	defer rows.Close()

	var overtimes []entity.Overtime
	for rows.Next() {
		var o entity.Overtime
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.PeriodID,
			&o.Date,
			&o.Hours,
			&o.CreatedAt,
			&o.UpdatedAt,
			&o.CreatedBy,
			&o.UpdatedBy,
		)
		if err != nil {
			r.Logger.Error(err, "OvertimeRepo - GetOvertimeByUserAndPeriod - Scan")
			return nil, err
		}
		overtimes = append(overtimes, o)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error(err, "OvertimeRepo - GetOvertimeByUserAndPeriod - rows.Err")
		return nil, err
	}

	return overtimes, nil
}

func (r *OvertimeRepo) GetTotalOvertimeByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (float64, error) {
	var total float64

	query := `
		SELECT COALESCE(SUM(hours), 0)
		FROM overtimes
		WHERE user_id = $1 AND date = $2
	`

	err := r.Pool.QueryRow(ctx, query, userID, date).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
