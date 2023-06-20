package mysql

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"

	apperr "github.com/rahmatrdn/go-skeleton/error"

	errwrap "github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	TrxSupportRepo
	GetByUserID(ctx context.Context, ID int64) (e []*entity.Wallet, err error)
	GetByID(ctx context.Context, ID int64) (e *entity.Wallet, err error)
	Create(ctx context.Context, dbTrx TrxObj, wallet *entity.Wallet, nonZeroVal bool) error
	LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Wallet, err error)
	Update(ctx context.Context, dbTrx TrxObj, params *entity.Wallet, changes *entity.Wallet) (err error)
	DeleteByID(ctx context.Context, dbTrx TrxObj, walletID int64) error
}

type Wallet struct {
	GormTrxSupport
}

func NewWalletRepository(mysql *config.Mysql) *Wallet {
	return &Wallet{GormTrxSupport{db: mysql.DB}}
}

func (u *Wallet) GetByUserID(ctx context.Context, userID int64) (e []*entity.Wallet, err error) {
	funcName := "WalletRepository.GetByUserID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = u.db.Where("user_id = ?", userID).Find(&e).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return e, err
}

func (u *Wallet) GetByID(ctx context.Context, ID int64) (e *entity.Wallet, err error) {
	funcName := "WalletRepository.GetByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = u.db.Where("id = ?", ID).Take(&e).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return e, err
}

func (u *Wallet) Create(ctx context.Context, dbTrx TrxObj, wallet *entity.Wallet, nonZeroVal bool) error {
	funcName := "WalletRepository.Create"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	cols := helper.NonZeroCols(wallet, nonZeroVal)
	return u.Trx(dbTrx).Select(cols).Create(&wallet).Error
}

func (u *Wallet) LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Wallet, err error) {
	funcName := "WalletRepository.LockByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = u.Trx(dbTrx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ID).Take(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (u *Wallet) Update(ctx context.Context, dbTrx TrxObj, params *entity.Wallet, changes *entity.Wallet) (err error) {
	funcName := "WalletRepository.Update"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	db := u.Trx(dbTrx).Model(params)
	if changes != nil {
		err = db.Updates(*changes).Error
	} else {
		err = db.Updates(helper.StructToMap(params, false)).Error
	}

	if err != nil {
		return errwrap.Wrap(err, funcName)
	}

	return nil
}

func (w *Wallet) DeleteByID(ctx context.Context, dbTrx TrxObj, walletID int64) error {
	funcName := "WalletRepository.DeleteByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	if err := w.Trx(dbTrx).Where("id = ?", walletID).Delete(&entity.Wallet{}).Error; err != nil {
		return err
	}

	return nil
}
