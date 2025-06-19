package overtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mock "github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/mocks"
)

type test struct {
	name string
	mock func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo)
	err  error
}

func TestSubmitOvertime(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	periodID := uuid.New()
	date := time.Now().Truncate(24 * time.Hour)

	tests := []test{
		{
			name: "success",
			mock: func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().
					GetActivePeriod(gomock.Any()).
					Return(&entity.AttendancePeriod{ID: periodID}, nil)

				repo.EXPECT().
					GetTotalOvertimeByUserAndDate(gomock.Any(), userID, date).
					Return(1.0, nil)

				repo.EXPECT().
					CreateOvertime(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "no active period",
			mock: func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().
					GetActivePeriod(gomock.Any()).
					Return(nil, errors.New("no active period"))
			},
			err: errors.New("422 - no active attendance period"),
		},
		{
			name: "overtime exceeds limit",
			mock: func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().
					GetActivePeriod(gomock.Any()).
					Return(&entity.AttendancePeriod{ID: periodID}, nil)

				repo.EXPECT().
					GetTotalOvertimeByUserAndDate(gomock.Any(), userID, date).
					Return(2.5, nil)
			},
			err: errors.New("409 - total overtime cannot exceed 3 hours per day"),
		},
		{
			name: "db error get total",
			mock: func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().
					GetActivePeriod(gomock.Any()).
					Return(&entity.AttendancePeriod{ID: periodID}, nil)

				repo.EXPECT().
					GetTotalOvertimeByUserAndDate(gomock.Any(), userID, date).
					Return(0.0, errors.New("db error"))
			},
			err: errors.New("500 - failed to validate existing overtime"),
		},
		{
			name: "db error create overtime",
			mock: func(repo *mock.MockOvertimeRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().
					GetActivePeriod(gomock.Any()).
					Return(&entity.AttendancePeriod{ID: periodID}, nil)

				repo.EXPECT().
					GetTotalOvertimeByUserAndDate(gomock.Any(), userID, date).
					Return(1.0, nil)

				repo.EXPECT().
					CreateOvertime(gomock.Any(), gomock.Any()).
					Return(errors.New("insert error"))
			},
			err: errors.New("insert error"),
		},
	}

	for _, tc := range tests {
		tc := tc // avoid loop variable capture
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			overtimeRepo := mock.NewMockOvertimeRepo(ctrl)
			periodRepo := mock.NewMockAttendancePeriodRepo(ctrl)
			log := logger.New("debug")

			uc := NewOvertimeUsecase(*log, overtimeRepo, periodRepo)

			tc.mock(overtimeRepo, periodRepo)

			req := request.OvertimeRequest{
				UserID:   userID,
				DateTime: date,
				Hours:    1.5,
				ClientIP: "127.0.0.1",
			}

			err := uc.SubmitOvertime(context.Background(), req)

			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
