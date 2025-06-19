package reimbursement

import (
	"context"
	"fmt"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/google/uuid"
)

type ReimbursementUsecase struct {
	logger               logger.Logger
	repo                 repo.ReimbursementRepo
	attendancePeriodRepo repo.AttendancePeriodRepo
}

// NewReimbursementUsecase -.
func NewReimbursementUsecase(
	l logger.Logger,
	r repo.ReimbursementRepo,
	apr repo.AttendancePeriodRepo,
) *ReimbursementUsecase {
	return &ReimbursementUsecase{
		logger:               l,
		repo:                 r,
		attendancePeriodRepo: apr,
	}
}

func (uc *ReimbursementUsecase) SubmitReimbursement(ctx context.Context, req request.ReimbursementRequest) error {
	// Get active period
	period, err := uc.attendancePeriodRepo.GetActivePeriod(ctx)
	if err != nil {
		return fmt.Errorf("422 - no active attendance period")
	}

	reimbursement := &entity.Reimbursement{
		ID:          uuid.New(),
		UserID:      req.UserID,
		PeriodID:    period.ID,
		Amount:      req.Amount,
		Description: req.Description,
		IPAddress:   req.ClientIP,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   req.UserID,
		UpdatedBy:   req.UserID,
	}

	if err := uc.repo.CreateReimbursement(ctx, reimbursement); err != nil {
		uc.logger.Error(err, "ReimbursementUsecase - SubmitReimbursement - uc.repo.CreateReimbursement")
		return err
	}

	return nil
}
