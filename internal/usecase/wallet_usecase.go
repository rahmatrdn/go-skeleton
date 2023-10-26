package usecase

import (
	"context"
	"fmt"
	"reflect"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
)

type Wallet struct {
	validatorUsecase ValidatorUsecase
	walletRepo       mysql.WalletRepository
	logger           LogUsecase
}

func NewWalletUsecase(
	validatorUsecase ValidatorUsecase,
	walletRepo mysql.WalletRepository,
	logger LogUsecase,
) *Wallet {
	return &Wallet{validatorUsecase, walletRepo, logger}
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
	captureFieldError := map[string]interface{}{"user_id": fmt.Sprint(userID)}

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
	captureFieldError := map[string]interface{}{"wallet_id": fmt.Sprint(walletID)}

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

type TableFields struct {
	StructFieldName string
	Field           string
}

func GetFieldsSlice(v interface{}) []TableFields {
	t := reflect.TypeOf(v)
	fieldsSlice := []TableFields{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldsSlice = append(fieldsSlice, TableFields{
			StructFieldName: field.Name,
			Field:           field.Tag.Get("json"),
		})
	}

	return fieldsSlice
}

func GetFieldByName(data []TableFields, structName string) *TableFields {
	for _, f := range data {
		if f.StructFieldName == structName {
			return &f
		}
	}
	return nil
}

func (w *Wallet) Create(ctx context.Context, walletReq *entity.WalletReq) (*entity.WalletResponse, error) {
	funcName := "WalletUsecase.Create"
	captureFieldError := map[string]interface{}{"user_name": fmt.Sprint(walletReq.UserName)}

	tableFields := GetFieldsSlice(entity.WalletReq{})

	wallet := &mentity.Wallet{
		UserName: walletReq.UserName,
		Balance:  walletReq.Balance,
	}

	errorData := w.validatorUsecase.Validate(walletReq)

	for _, v := range errorData {
		// Get nama kolom yang ada pada table di database
		tableField := GetFieldByName(tableFields, v.FailedField)

		helper.Dump(v)
		helper.Dump(tableField.Field)
	}

	// Struct Validation is like form validation
	if errMsg := w.validatorUsecase.ValidateWithMessage(walletReq); errMsg != "null" {
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

	captureFieldError := map[string]interface{}{"wallet_id": fmt.Sprint(walletID)}

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
	captureFieldError := map[string]interface{}{"wallet_id": fmt.Sprint(walletID)}

	err := w.walletRepo.DeleteByID(ctx, nil, walletID)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
