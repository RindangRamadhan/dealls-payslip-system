package attendance

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

type AttendanceUsecase struct {
	logger               logger.Logger
	repo                 repo.AttendanceRepo
	attendancePeriodRepo repo.AttendancePeriodRepo
}

// NewAttendanceUsecase -.
func NewAttendanceUsecase(
	l logger.Logger,
	r repo.AttendanceRepo,
	apr repo.AttendancePeriodRepo,
) *AttendanceUsecase {
	return &AttendanceUsecase{
		logger:               l,
		repo:                 r,
		attendancePeriodRepo: apr,
	}
}

func (uc *AttendanceUsecase) SubmitAttendance(ctx context.Context, req request.AttendanceRequest) error {
	// 1. Retrieve the currently active attendance period
	period, err := uc.attendancePeriodRepo.GetActivePeriod(ctx)
	if err != nil {
		return fmt.Errorf("422 - no active attendance period")
	}

	// 2. Check if there's already an attendance record for this date
	existing, err := uc.repo.GetAttendanceByDate(ctx, req.UserID, period.ID, req.DateTime)
	if err != nil {
		uc.logger.Error(err, "AttendanceUsecase - SubmitAttendance - GetAttendanceByDate")
		return err
	}

	switch req.Type {
	case "check_in":
		// Prevent duplicate check-in on the same date
		if existing != nil && !existing.CheckIn.IsZero() {
			return fmt.Errorf("409 - attendance already submitted for this date")
		}

		// Create new attendance record with check-in
		attendance := &entity.Attendance{
			ID:        uuid.New(),
			UserID:    req.UserID,
			PeriodID:  period.ID,
			Date:      req.DateTime,
			CheckIn:   req.CheckInTime,
			CheckOut:  req.CheckOutTime,
			IPAddress: req.ClientIP,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: req.UserID,
			UpdatedBy: req.UserID,
		}

		if err := uc.repo.CreateAttendance(ctx, attendance); err != nil {
			uc.logger.Error(err, "SubmitAttendance - CreateAttendance")
			return err
		}

	case "check_out":
		// Ensure a check-in exists before check-out
		if existing == nil {
			return fmt.Errorf("404 - check-in record not found for this date")
		}

		// Prevent duplicate check-out
		if existing.CheckOut != nil {
			return fmt.Errorf("409 - already checked out for this date")
		}

		if err := uc.repo.UpdateCheckOut(ctx, existing.ID, *req.CheckOutTime, req.UserID); err != nil {
			uc.logger.Error(err, "SubmitAttendance - UpdateCheckOut")
			return err
		}

	default:
		return fmt.Errorf("400 - invalid attendance type: must be check_in or check_out")
	}

	return nil
}
