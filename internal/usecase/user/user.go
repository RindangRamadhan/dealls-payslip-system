package user

import (
	"context"
	"fmt"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repo.UserRepo
}

// NewUserUsecase -.
func NewUserUsecase(r repo.UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: r,
	}
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserUsecase - GetUserByID - s.repo.GetUserByID: %w", err)
	}

	return user, nil
}

func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("UserUsecase - GetUserByUsername - s.repo.GetGetUserByUsername: %w", err)
	}

	return user, nil
}
