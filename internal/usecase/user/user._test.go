package user

import (
	"context"
	"errors"
	"testing"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mock "github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/mocks"
)

type testCase struct {
	name      string
	mock      func(repo *mock.MockUserRepo)
	exec      func(uc *UserUsecase) (*entity.User, error)
	expectErr error
}

func TestUserUsecase(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	username := "john_doe"
	exampleUser := &entity.User{
		ID:       userID,
		Username: username,
	}

	tests := []testCase{
		{
			name: "GetUserByID success",
			mock: func(repo *mock.MockUserRepo) {
				repo.EXPECT().GetUserByID(gomock.Any(), userID).Return(exampleUser, nil)
			},
			exec: func(uc *UserUsecase) (*entity.User, error) {
				return uc.GetUserByID(context.Background(), userID)
			},
			expectErr: nil,
		},
		{
			name: "GetUserByID not found",
			mock: func(repo *mock.MockUserRepo) {
				repo.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.New("user not found"))
			},
			exec: func(uc *UserUsecase) (*entity.User, error) {
				return uc.GetUserByID(context.Background(), userID)
			},
			expectErr: errors.New("user not found"),
		},
		{
			name: "GetUserByUsername success",
			mock: func(repo *mock.MockUserRepo) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(exampleUser, nil)
			},
			exec: func(uc *UserUsecase) (*entity.User, error) {
				return uc.GetUserByUsername(context.Background(), username)
			},
			expectErr: nil,
		},
		{
			name: "GetUserByUsername not found",
			mock: func(repo *mock.MockUserRepo) {
				repo.EXPECT().GetUserByUsername(gomock.Any(), username).Return(nil, errors.New("user not found"))
			},
			exec: func(uc *UserUsecase) (*entity.User, error) {
				return uc.GetUserByUsername(context.Background(), username)
			},
			expectErr: errors.New("user not found"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mock.NewMockUserRepo(ctrl)
			tc.mock(repo)

			uc := NewUserUsecase(repo)
			result, err := tc.exec(uc)

			if tc.expectErr != nil {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, exampleUser.ID, result.ID)
			}
		})
	}
}
