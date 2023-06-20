package usecase

import (
	"context"
	"fmt"

	errwrap "github.com/pkg/errors"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"
	apperr "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/error"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mysql"
	mentity "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mysql/entity"
)

type Wallet struct {
	walletRepo mysql.WalletRepository
	logger     LogUsecase
}

func NewWalletUsecase(
	walletRepo mysql.WalletRepository,
	logger LogUsecase,
) *Wallet {
	return &Wallet{walletRepo, logger}
}

type WalletUsecase interface {
	GetByUserID(ctx context.Context, userID int64) (res []*entity.WalletResponse, err error)
	GetByID(ctx context.Context, walletID int64) (*entity.WalletResponse, error)
	Create(ctx context.Context, walletReq *entity.WalletReq) (*entity.WalletResponse, error)
	UpdateByID(ctx context.Context, walletReq entity.WalletReq) error
	DeleteByID(ctx context.Context, walletID int64) error
}

func (w *Wallet) GetByUserID(ctx context.Context, userID int64) (res []*entity.WalletResponse, err error) {
	funcName := "WalletUsecase.GetByUserID"
	captureFieldError := map[string]string{"user_id": fmt.Sprint(userID)}

	result, err := w.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.GetByUserID", funcName, err, captureFieldError, "")

		return nil, err
	}

	for _, v := range result {
		res = append(res, &entity.WalletResponse{
			ID:       v.ID,
			UserName: v.UserName,
			Balance:  v.Balance,
		})
	}

	return res, nil
}

func (w *Wallet) GetByID(ctx context.Context, walletID int64) (*entity.WalletResponse, error) {
	funcName := "WalletUsecase.GetByID"
	captureFieldError := map[string]string{"wallet_id": fmt.Sprint(walletID)}

	data, err := w.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.WalletResponse{
		ID:       data.ID,
		UserName: data.UserName,
		Balance:  data.Balance,
	}, nil
}

func (w *Wallet) Create(ctx context.Context, walletReq *entity.WalletReq) (*entity.WalletResponse, error) {
	funcName := "WalletUsecase.Create"
	captureFieldError := map[string]string{"user_name": fmt.Sprint(walletReq.UserName)}

	wallet := &mentity.Wallet{
		UserName: walletReq.UserName,
		Balance:  walletReq.Balance,
	}

	// Struct Validation is like form validation
	if errMsg := ValidateStruct(walletReq); errMsg != "null" {
		fmt.Println(errMsg)
		return nil, errwrap.Wrap(fmt.Errorf(apperr.INVALID_PAYLOAD_CODE), errMsg)
	}

	err := w.walletRepo.Create(ctx, nil, wallet, false)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.WalletResponse{
		ID:       wallet.ID,
		UserName: wallet.UserName,
		Balance:  wallet.Balance,
	}, nil
}

func (w *Wallet) UpdateByID(ctx context.Context, walletReq entity.WalletReq) error {
	funcName := "WalletUsecase.UpdateByID"
	walletID := walletReq.ID

	captureFieldError := map[string]string{"wallet_id": fmt.Sprint(walletID)}

	if err := mysql.DBTransaction(w.walletRepo, func(trx mysql.TrxObj) error {
		lockedWallet, err := w.walletRepo.LockByID(ctx, trx, walletID)
		if err != nil {
			w.logger.Log(entity.LogError, "walletRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}

		if err := w.walletRepo.Update(ctx, trx, lockedWallet, &mentity.Wallet{
			UserName: walletReq.UserName,
			Balance:  walletReq.Balance,
		}); err != nil {
			w.logger.Log(entity.LogError, "walletRepo.Update", funcName, err, captureFieldError, "")

			return err
		}

		return nil
	}); err != nil {
		w.logger.Log(entity.LogError, "walletRepo.DBTransaction", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}

func (w *Wallet) DeleteByID(ctx context.Context, walletID int64) error {
	funcName := "WalletUsecase.UpdateByID"
	captureFieldError := map[string]string{"wallet_id": fmt.Sprint(walletID)}

	err := w.walletRepo.DeleteByID(ctx, nil, walletID)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
