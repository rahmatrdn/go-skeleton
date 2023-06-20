package usecase_test

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/usecase"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/tests/fixture/factory"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/tests/mocks"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WalletUsecaseTestSuite struct {
	suite.Suite

	walletRepo *mocks.WalletRepository
	logger     *mocks.LogUsecase
	usecase    usecase.WalletUsecase
	trxObj     *mocks.TrxObj
}

func (s *WalletUsecaseTestSuite) SetupTest() {
	s.walletRepo = &mocks.WalletRepository{}
	s.logger = &mocks.LogUsecase{}

	s.usecase = usecase.NewWalletUsecase(s.walletRepo, s.logger)
	s.trxObj = &mocks.TrxObj{}

	s.walletRepo.On("Begin").Return(s.trxObj, nil)
	s.trxObj.On("Rollback", mock.Anything).Return(nil)
	s.trxObj.On("Commit", mock.Anything).Return(nil)
}

func TestWalletUsecase(t *testing.T) {
	suite.Run(t, new(WalletUsecaseTestSuite))
}

func (s *WalletUsecaseTestSuite) TestGetByID() {
	wallet := factory.StubbedWallet()

	ID := int64(1)

	ctx := context.Background()

	testcases := []struct {
		name     string
		mockFunc func()
		want     *entity.WalletResponse
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.walletRepo.On("GetByID", ctx, ID).Return(wallet, nil).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			want: &entity.WalletResponse{
				ID:       wallet.ID,
				UserName: wallet.UserName,
				Balance:  wallet.Balance,
			},
		},
		{
			name: "failure",
			mockFunc: func() {
				s.walletRepo.On("GetByID", ctx, ID).Return(nil, fmt.Errorf("ERROR")).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			got, err := s.usecase.GetByID(ctx, ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *WalletUsecaseTestSuite) TestCreate() {
	wallet := factory.StubbedWallet()

	params := &entity.WalletReq{
		UserName: wallet.UserName,
		Balance:  wallet.Balance,
		UserID:   rand.Int63(),
	}
	ctx := context.Background()

	testCases := []struct {
		name     string
		params   *entity.WalletReq
		mockFunc func(params *entity.WalletReq)
		want     *entity.WalletResponse
		wantErr  bool
	}{
		{
			name:   "success",
			params: params,
			mockFunc: func(params *entity.WalletReq) {
				s.walletRepo.On("Create", ctx, mock.Anything, mock.Anything, false).Return(nil).Once()
			},
			wantErr: false,
			want: &entity.WalletResponse{
				ID:       0,
				UserName: wallet.UserName,
				Balance:  wallet.Balance,
			},
		},
		{
			name:   "fail create to repo",
			params: params,
			mockFunc: func(params *entity.WalletReq) {
				s.walletRepo.On("Create", ctx, mock.Anything, mock.Anything, false).Return(fmt.Errorf("ERROR")).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			wantErr: true,
			want:    nil,
		},
		{
			name:   "fail struct validation",
			params: params,
			mockFunc: func(params *entity.WalletReq) {
				params.UserName = ""
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.params)

			got, err := s.usecase.Create(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *WalletUsecaseTestSuite) TestUpdateByID() {
	wallet := factory.StubbedWallet()
	ctx := context.Background()
	ID := int64(1)

	params := entity.WalletReq{
		ID:       1,
		UserName: wallet.UserName,
		Balance:  wallet.Balance,
	}

	testCases := []struct {
		name     string
		params   entity.WalletReq
		mockFunc func(params entity.WalletReq)
		wantErr  bool
	}{
		{
			name:   "success",
			params: params,
			mockFunc: func(params entity.WalletReq) {
				s.walletRepo.On("LockByID", ctx, mock.Anything, ID).Return(wallet, nil).Once()
				s.walletRepo.On("Update", ctx, mock.Anything, wallet, mock.Anything).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "fail lock by id",
			params: params,
			mockFunc: func(params entity.WalletReq) {
				s.walletRepo.On("LockByID", ctx, mock.Anything, ID).Return(nil, fmt.Errorf("Error")).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			wantErr: true,
		},
		{
			name:   "fail update in repo",
			params: params,
			mockFunc: func(params entity.WalletReq) {
				s.walletRepo.On("LockByID", ctx, mock.Anything, ID).Return(wallet, nil).Once()
				s.walletRepo.On("Update", ctx, mock.Anything, wallet, mock.Anything).Return(fmt.Errorf("Error")).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.params)

			err := s.usecase.UpdateByID(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (s *WalletUsecaseTestSuite) TestDeleteByID() {
	wallet := factory.StubbedWallet()
	ctx := context.Background()

	testCases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.walletRepo.On("DeleteByID", ctx, nil, wallet.ID).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "fail",
			mockFunc: func() {
				s.walletRepo.On("DeleteByID", ctx, nil, wallet.ID).Return(fmt.Errorf("ERROR")).Once()
				s.logger.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wallet, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.usecase.DeleteByID(ctx, wallet.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAddressByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
