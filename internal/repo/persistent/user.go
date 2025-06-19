package persistent

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"
	"github.com/google/uuid"
)

type UserRepo struct {
	*postgres.Postgres
}

// UserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetAllEmployees(ctx context.Context) ([]entity.User, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"username",
			"password",
			"is_admin",
			"salary",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("users").
		Where(squirrel.Eq{"is_admin": false}).
		OrderBy("username").
		ToSql()
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetAllEmployees - ToSql")
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetAllEmployees - Query")
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.IsAdmin,
			&u.Salary,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.CreatedBy,
			&u.UpdatedBy,
		)
		if err != nil {
			r.Logger.Error(err, "UserRepo - GetAllEmployees - Scan")
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		r.Logger.Error(err, "UserRepo - GetAllEmployees - rows.Err")
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"username",
			"password",
			"is_admin",
			"salary",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetUserByUsername - ToSql")
		return nil, err
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.Salary,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetUserByUsername - QueryRow")
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select(
			"id",
			"username",
			"password",
			"is_admin",
			"salary",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetUserByID - ToSql")
		return nil, err
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.Salary,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)
	if err != nil {
		r.Logger.Error(err, "UserRepo - GetUserByID - QueryRow")
		return nil, err
	}

	return &user, nil
}
