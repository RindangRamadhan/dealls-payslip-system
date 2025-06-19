package attendanceperiod

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mock "github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/mocks"
)

type test struct {
	name string
	mock func(repo *mock.MockAttendancePeriodRepo)
	err  error
}

func TestCreateAttendancePeriod(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Second)
	startDate := now
	endDate := now.AddDate(0, 0, 7) // 7 hari ke depan

	period := &entity.AttendancePeriod{
		ID:        uuid.New(),
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []test{
		{
			name: "success",
			mock: func(repo *mock.MockAttendancePeriodRepo) {
				repo.EXPECT().CreateAttendancePeriod(gomock.Any(), period).Return(nil)
			},
			err: nil,
		},
		{
			name: "failed to create attendance period",
			mock: func(repo *mock.MockAttendancePeriodRepo) {
				repo.EXPECT().CreateAttendancePeriod(gomock.Any(), period).Return(errors.New("db error"))
			},
			err: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		tc := tc // avoid variable capture
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mock.NewMockAttendancePeriodRepo(ctrl)
			log := logger.New("debug")

			uc := NewAttendancePeriodUsecase(*log, repo)

			// Setup mock expectation
			tc.mock(repo)

			err := uc.CreateAttendancePeriod(context.Background(), period)

			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
