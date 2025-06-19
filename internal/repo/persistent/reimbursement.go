package persistent

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type ReimbursementRepo struct {
	*postgres.Postgres
}

// ReimbursementRepo -.
func NewReimbursementRepo(pg *postgres.Postgres) *ReimbursementRepo {
	return &ReimbursementRepo{pg}
}

func (r *ReimbursementRepo) CreateReimbursement(ctx context.Context, reimbursement *entity.Reimbursement) error {
	sql, args, err := r.Builder.
		Insert("reimbursements").
		Columns(
			"id",
			"user_id",
			"period_id",
			"amount",
			"description",
			"ip_address",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Values(
			reimbursement.ID,
			reimbursement.UserID,
			reimbursement.PeriodID,
			reimbursement.Amount,
			reimbursement.Description,
			reimbursement.IPAddress,
			reimbursement.CreatedAt,
			reimbursement.UpdatedAt,
			reimbursement.CreatedBy,
			reimbursement.UpdatedBy,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "ReimbursementRepo - CreateReimbursement - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "ReimbursementRepo - CreateReimbursement - Exec")
		return err
	}

	return nil
}

func (r *ReimbursementRepo) GetReimbursementsByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]entity.Reimbursement, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"user_id",
			"period_id",
			"amount",
			"description",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("reimbursements").
		Where(squirrel.Eq{"user_id": userID, "period_id": periodID}).
		OrderBy("created_at").
		ToSql()
	if err != nil {
		r.Logger.Error(err, "ReimbursementRepo - GetReimbursementsByUserAndPeriod - ToSql")
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "ReimbursementRepo - GetReimbursementsByUserAndPeriod - Query")
		return nil, err
	}
	defer rows.Close()

	var reimbursements []entity.Reimbursement
	for rows.Next() {
		var rs entity.Reimbursement
		err := rows.Scan(
			&rs.ID,
			&rs.UserID,
			&rs.PeriodID,
			&rs.Amount,
			&rs.Description,
			&rs.CreatedAt,
			&rs.UpdatedAt,
			&rs.CreatedBy,
			&rs.UpdatedBy,
		)
		if err != nil {
			r.Logger.Error(err, "ReimbursementRepo - GetReimbursementsByUserAndPeriod - Scan")
			return nil, err
		}
		reimbursements = append(reimbursements, rs)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error(err, "ReimbursementRepo - GetReimbursementsByUserAndPeriod - rows.Err")
		return nil, err
	}

	return reimbursements, nil
}
