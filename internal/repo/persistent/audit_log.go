package persistent

import (
	"context"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
)

type AuditLogRepo struct {
	*postgres.Postgres
}

// AuditLogRepo -.
func NewAuditLogRepo(pg *postgres.Postgres) *AuditLogRepo {
	return &AuditLogRepo{pg}
}

func (r *AuditLogRepo) CreateAuditLog(ctx context.Context, log *entity.AuditLog) error {
	sql, args, err := r.Builder.
		Insert("audit_logs").
		Columns(
			"id",
			"user_id",
			"action",
			"table_name",
			"record_id",
			"old_values",
			"new_values",
			"ip_address",
			"request_id",
			"created_at",
		).
		Values(
			log.ID,
			log.UserID,
			log.Action,
			log.TableName,
			log.RecordID,
			log.OldValues,
			log.NewValues,
			log.IPAddress,
			log.RequestID,
			log.CreatedAt,
		).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "AuditLogRepo - CreateAuditLog - ToSql")
		return err
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "AuditLogRepo - CreateAuditLog - Exec")
		return err
	}

	return nil
}
