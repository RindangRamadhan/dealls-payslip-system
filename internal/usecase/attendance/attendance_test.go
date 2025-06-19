package attendance

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

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func(repo *mock.MockAttendanceRepo, periodRepo *mock.MockAttendancePeriodRepo)
	err  error
}

func TestSubmitAttendance(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	periodID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)
	checkIn := time.Date(today.Year(), today.Month(), today.Day(), 9, 0, 0, 0, today.Location())
	checkOut := time.Date(today.Year(), today.Month(), today.Day(), 17, 0, 0, 0, today.Location())

	tests := []test{
		{
			name: "check-in success",
			mock: func(repo *mock.MockAttendanceRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				repo.EXPECT().GetAttendanceByDate(gomock.Any(), userID, periodID, today).Return(nil, nil)
				repo.EXPECT().CreateAttendance(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "check-in already exists",
			mock: func(repo *mock.MockAttendanceRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				existing := &entity.Attendance{CheckIn: &checkIn}
				repo.EXPECT().GetAttendanceByDate(gomock.Any(), userID, periodID, today).Return(existing, nil)
			},
			err: errors.New("409 - attendance already submitted for this date"),
		},
		{
			name: "check-out success",
			mock: func(repo *mock.MockAttendanceRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				existing := &entity.Attendance{ID: uuid.New(), CheckIn: &checkIn}
				repo.EXPECT().GetAttendanceByDate(gomock.Any(), userID, periodID, today).Return(existing, nil)
				repo.EXPECT().UpdateCheckOut(gomock.Any(), existing.ID, checkOut, userID).Return(nil)
			},
			err: nil,
		},
		{
			name: "check-out without check-in",
			mock: func(repo *mock.MockAttendanceRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				repo.EXPECT().GetAttendanceByDate(gomock.Any(), userID, periodID, today).Return(nil, nil)
			},
			err: errors.New("404 - check-in record not found for this date"),
		},
	}

	for _, tc := range tests {
		tc := tc // avoid variable capture
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mock.NewMockAttendanceRepo(ctrl)
			periodRepo := mock.NewMockAttendancePeriodRepo(ctrl)
			log := logger.New("debug")

			uc := NewAttendanceUsecase(*log, repo, periodRepo)

			// Setup mock expectation
			tc.mock(repo, periodRepo)

			req := request.AttendanceRequest{
				UserID:       userID,
				DateTime:     today,
				CheckInTime:  &checkIn,
				CheckOutTime: &checkOut,
				ClientIP:     "127.0.0.1",
				Type:         "check_in",
			}

			if tc.name == "check-out success" || tc.name == "check-out without check-in" {
				req.Type = "check_out"
			}
			if tc.name == "invalid type" {
				req.Type = "invalid"
			}

			err := uc.SubmitAttendance(context.Background(), req)

			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
