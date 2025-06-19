package persistent

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type AttendanceRepo struct {
	*postgres.Postgres
}

// AttendanceRepo -.
func NewAttendanceRepo(pg *postgres.Postgres) *AttendanceRepo {
	return &AttendanceRepo{pg}
}

func (r *AttendanceRepo) CreateAttendance(ctx context.Context, attendance *entity.Attendance) error {
	sql, args, err := r.Builder.
		Insert("attendances").
		Columns(
			"id",
			"user_id",
			"period_id",
			"date",
			"check_in",
			"check_out",
			"ip_address",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			attendance.ID,
			attendance.UserID,
			attendance.PeriodID,
			attendance.Date,
			attendance.CheckIn,
			attendance.CheckOut,
			attendance.IPAddress,
			attendance.CreatedAt,
			attendance.UpdatedAt,
			attendance.CreatedBy,
			attendance.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - CreateAttendance - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - CreateAttendance - Exec")
		return err
	}

	return nil
}

func (r *AttendanceRepo) UpdateCheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time, updatedBy uuid.UUID) error {
	sql, args, err := r.Builder.
		Update("attendances").
		Set("check_out", checkOutTime).
		Set("updated_at", time.Now()).
		Set("updated_by", updatedBy).
		Where(squirrel.Eq{"id": attendanceID}).
		ToSql()

	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - UpdateCheckOut - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - UpdateCheckOut - Exec")
		return err
	}

	return nil
}

func (r *AttendanceRepo) GetAttendanceByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Attendance, error) {
	sql := `
		SELECT DISTINCT ON (date)
			id,
			user_id,
			period_id,
			date,
			check_in,
			check_out,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM attendances
		WHERE user_id = $1 AND period_id = $2
		ORDER BY date, created_at DESC
	`

	rows, err := r.Pool.Query(ctx, sql, userID, periodID)
	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - GetAttendanceByUserAndPeriod - Query")
		return nil, err
	}
	defer rows.Close()

	var attendances []entity.Attendance
	for rows.Next() {
		var a entity.Attendance
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.PeriodID,
			&a.Date,
			&a.CheckIn,
			&a.CheckOut,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.CreatedBy,
			&a.UpdatedBy,
		)
		if err != nil {
			r.Logger.Error(err, "AttendanceRepo - GetAttendanceByUserAndPeriod - Scan")
			return nil, err
		}
		attendances = append(attendances, a)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error(err, "AttendanceRepo - GetAttendanceByUserAndPeriod - rows.Err")
		return nil, err
	}

	return attendances, nil
}

func (r *AttendanceRepo) GetAttendanceByDate(ctx context.Context, userID, periodID uuid.UUID, date time.Time) (*entity.Attendance, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"user_id",
			"period_id",
			"date",
			"check_in",
			"check_out",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("attendances").
		Where(squirrel.Eq{
			"user_id":   userID,
			"period_id": periodID,
			"date":      date,
		}).
		Limit(1).
		ToSql()

	if err != nil {
		r.Logger.Error(err, "AttendanceRepo - GetAttendanceByDate - ToSql")
		return nil, err
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	var a entity.Attendance
	err = row.Scan(
		&a.ID,
		&a.UserID,
		&a.PeriodID,
		&a.Date,
		&a.CheckIn,
		&a.CheckOut,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.CreatedBy,
		&a.UpdatedBy,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}

		r.Logger.Error(err, "AttendanceRepo - GetAttendanceByDate - Scan")
		return nil, err
	}

	return &a, nil
}
