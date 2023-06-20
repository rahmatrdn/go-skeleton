package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TrxSupportRepo interface {
	Begin() (TrxObj, error)
}

type TrxObj interface {
	Commit() error
	Rollback() error
}

// GormTrxSupport parent mysqlrepo
type GormTrxSupport struct {
	db *gorm.DB
}

type GormTrxObj struct {
	db *gorm.DB
}

// Begin Begin db transaction
func (repo *GormTrxSupport) Begin() (TrxObj, error) {
	trx := repo.db.Begin(&sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	return &GormTrxObj{db: trx}, trx.Error
}

// Trx get transaction
func (repo *GormTrxSupport) Trx(trx TrxObj) *gorm.DB {
	gormTrx, ok := trx.(*GormTrxObj)
	if ok {
		return gormTrx.db
	}

	return repo.db
}

// Commit Commit db transaction
func (trx *GormTrxObj) Commit() error {
	return trx.db.Commit().Error
}

// Rollback rollback trx
func (trx *GormTrxObj) Rollback() error {
	return trx.db.Rollback().Error
}

// DBTransaction usecase with db transaction
func DBTransaction(repo TrxSupportRepo, callback func(TrxObj) error) (err error) {
	functionName := "DBTransaction"
	commit := false
	trx, err := repo.Begin()
	if err != nil {
		return err
	}

	defer func(commit *bool, repo TrxSupportRepo, trx TrxObj) {
		if !*commit {
			if rErr := trx.Rollback(); rErr != nil {
				if err == nil {
					err = rErr
				} else {
					err = errors.Wrap(rErr, err.Error())
				}
			}
		}
	}(&commit, repo, trx)

	// call the callback function
	err = callback(trx)
	if err != nil {
		return err
	}

	if err = trx.Commit(); err != nil {
		return errors.Wrap(err, functionName)
	}
	commit = true

	return err
}
