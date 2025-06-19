package reimbursement

import (
	"context"
	"errors"
	"testing"

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
	mock func(repo *mock.MockReimbursementRepo, periodRepo *mock.MockAttendancePeriodRepo)
	err  error
}

func TestSubmitReimbursement(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	periodID := uuid.New()

	tests := []test{
		{
			name: "success",
			mock: func(repo *mock.MockReimbursementRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				repo.EXPECT().CreateReimbursement(gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "no active attendance period",
			mock: func(repo *mock.MockReimbursementRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(nil, errors.New("not found"))
			},
			err: errors.New("422 - no active attendance period"),
		},
		{
			name: "create reimbursement failed",
			mock: func(repo *mock.MockReimbursementRepo, periodRepo *mock.MockAttendancePeriodRepo) {
				periodRepo.EXPECT().GetActivePeriod(gomock.Any()).Return(&entity.AttendancePeriod{ID: periodID}, nil)
				repo.EXPECT().CreateReimbursement(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			err: errors.New("db error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mock.NewMockReimbursementRepo(ctrl)
			periodRepo := mock.NewMockAttendancePeriodRepo(ctrl)
			log := logger.New("debug")

			uc := NewReimbursementUsecase(*log, repo, periodRepo)

			tc.mock(repo, periodRepo)

			req := request.ReimbursementRequest{
				UserID:      userID,
				Amount:      100000,
				Description: "Transport",
				ClientIP:    "127.0.0.1",
			}

			err := uc.SubmitReimbursement(context.Background(), req)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
