package attendanceperiod

import (
	"context"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
)

type AttendancePeriodUsecase struct {
	logger logger.Logger
	repo   repo.AttendancePeriodRepo
}

// NewAttendancePeriodUsecase -.
func NewAttendancePeriodUsecase(l logger.Logger, r repo.AttendancePeriodRepo) *AttendancePeriodUsecase {
	return &AttendancePeriodUsecase{
		logger: l,
		repo:   r,
	}
}

func (uc *AttendancePeriodUsecase) CreateAttendancePeriod(ctx context.Context, req *entity.AttendancePeriod) error {
	err := uc.repo.CreateAttendancePeriod(ctx, req)
	if err != nil {
		uc.logger.Error(err, "AttendancePeriodUsecase - CreateAttendancePeriod - s.repo.CreateAttendancePeriod")
		return err
	}

	return nil
}
