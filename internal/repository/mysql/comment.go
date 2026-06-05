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

type ICommentRepository interface {
	TrxSupportRepo
	GetByPostID(ctx context.Context, postID int64) (result []*entity.Comment, err error)
	GetByID(ctx context.Context, ID int64) (result *entity.Comment, err error)
	Create(ctx context.Context, dbTrx TrxObj, params *entity.Comment, nonZeroVal bool) error
	LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Comment, err error)
	Update(ctx context.Context, dbTrx TrxObj, params *entity.Comment, changes *entity.Comment) (err error)
	DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error
}

type CommentRepository struct {
	GormTrxSupport
}

func NewCommentRepository(mysql *config.Mysql) *CommentRepository {
	return &CommentRepository{GormTrxSupport{db: mysql.DB}}
}

func (r *CommentRepository) GetByPostID(ctx context.Context, postID int64) (result []*entity.Comment, err error) {
	funcName := "CommentRepository.GetByPostID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM comments WHERE post_id = ? ORDER BY created_at ASC", postID).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *CommentRepository) GetByID(ctx context.Context, ID int64) (result *entity.Comment, err error) {
	funcName := "CommentRepository.GetByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM comments WHERE id = ? LIMIT 1", ID).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *CommentRepository) Create(ctx context.Context, dbTrx TrxObj, params *entity.Comment, nonZeroVal bool) error {
	funcName := "CommentRepository.Create"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	cols := helper.NonZeroCols(params, nonZeroVal)
	return r.Trx(dbTrx).Select(cols).Create(&params).Error
}

func (r *CommentRepository) LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.Comment, err error) {
	funcName := "CommentRepository.LockByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.Trx(dbTrx).
		Raw("SELECT * FROM comments WHERE id = ? FOR UPDATE", ID).
		Scan(&result).Error

	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *CommentRepository) Update(ctx context.Context, dbTrx TrxObj, params *entity.Comment, changes *entity.Comment) (err error) {
	funcName := "CommentRepository.Update"

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

func (r *CommentRepository) DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error {
	funcName := "CommentRepository.DeleteByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	return r.Trx(dbTrx).
		Exec("DELETE FROM comments WHERE id = ?", id).
		Error
}
