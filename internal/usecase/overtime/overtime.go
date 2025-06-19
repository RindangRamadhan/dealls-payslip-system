package overtime

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

type OvertimeUsecase struct {
	logger               logger.Logger
	repo                 repo.OvertimeRepo
	attendancePeriodRepo repo.AttendancePeriodRepo
}

// NewOvertimeUsecase -.
func NewOvertimeUsecase(
	l logger.Logger,
	r repo.OvertimeRepo,
	apr repo.AttendancePeriodRepo,
) *OvertimeUsecase {
	return &OvertimeUsecase{
		logger:               l,
		repo:                 r,
		attendancePeriodRepo: apr,
	}
}

func (uc *OvertimeUsecase) SubmitOvertime(ctx context.Context, req request.OvertimeRequest) error {
	// Get active period
	period, err := uc.attendancePeriodRepo.GetActivePeriod(ctx)
	if err != nil {
		return fmt.Errorf("422 - no active attendance period")
	}

	existingHours, err := uc.repo.GetTotalOvertimeByUserAndDate(ctx, req.UserID, req.DateTime)
	if err != nil {
		uc.logger.Error(err, "OvertimeUsecase - SubmitOvertime - GetTotalOvertimeByUserAndDate")
		return fmt.Errorf("500 - failed to validate existing overtime")
	}

	totalHours := existingHours + req.Hours
	if totalHours > 3 {
		return fmt.Errorf("409 - total overtime cannot exceed 3 hours per day")
	}

	overtime := &entity.Overtime{
		ID:        uuid.New(),
		UserID:    req.UserID,
		PeriodID:  period.ID,
		Date:      req.DateTime,
		Hours:     req.Hours,
		IPAddress: req.ClientIP,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
	}

	if err := uc.repo.CreateOvertime(ctx, overtime); err != nil {
		uc.logger.Error(err, "OvertimeUsecase - SubmitAttendance - uc.repo.CreateOvertime")
		return err
	}

	return nil
}
