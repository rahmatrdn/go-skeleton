package mysql

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"

	apperr "github.com/rahmatrdn/go-skeleton/error"

	errwrap "github.com/pkg/errors"
	"gorm.io/gorm"
)

type IPostRepository interface {
	TrxSupportRepo
	GetAll(ctx context.Context) (result []*entity.Post, err error)
	GetByID(ctx context.Context, ID int64) (result *entity.Post, err error)
	GetBySlug(ctx context.Context, slug string) (result *entity.Post, err error)
	Create(ctx context.Context, dbTrx TrxObj, params *entity.Post, nonZeroVal bool) error
	LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Post, err error)
	Update(ctx context.Context, dbTrx TrxObj, params *entity.Post, changes *entity.Post) (err error)
	DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error
}

type PostRepository struct {
	GormTrxSupport
}

func NewPostRepository(mysql *config.Mysql) *PostRepository {
	return &PostRepository{GormTrxSupport{db: mysql.DB}}
}

func (r *PostRepository) GetAll(ctx context.Context) (result []*entity.Post, err error) {
	funcName := "PostRepository.GetAll"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM posts ORDER BY created_at DESC").Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *PostRepository) GetByID(ctx context.Context, ID int64) (result *entity.Post, err error) {
	funcName := "PostRepository.GetByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM posts WHERE id = ? LIMIT 1", ID).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *PostRepository) GetBySlug(ctx context.Context, slug string) (result *entity.Post, err error) {
	funcName := "PostRepository.GetBySlug"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM posts WHERE slug = ? LIMIT 1", slug).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *PostRepository) Create(ctx context.Context, dbTrx TrxObj, params *entity.Post, nonZeroVal bool) error {
	funcName := "PostRepository.Create"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	cols := helper.NonZeroCols(params, nonZeroVal)
	return r.Trx(dbTrx).Select(cols).Create(&params).Error
}

func (r *PostRepository) LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Post, err error) {
	funcName := "PostRepository.LockByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.Trx(dbTrx).
		Raw("SELECT * FROM posts WHERE id = ? FOR UPDATE", ID).
		Scan(&result).Error

	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *PostRepository) Update(ctx context.Context, dbTrx TrxObj, params *entity.Post, changes *entity.Post) (err error) {
	funcName := "PostRepository.Update"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	db := r.Trx(dbTrx).Model(params)
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

func (r *PostRepository) DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error {
	funcName := "PostRepository.DeleteByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	return r.Trx(dbTrx).
		Exec("DELETE FROM posts WHERE id = ?", id).
		Error
}
