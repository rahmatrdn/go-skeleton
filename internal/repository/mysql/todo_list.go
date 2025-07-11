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

type ITodoListRepository interface {
	TrxSupportRepo
	GetByUserID(ctx context.Context, ID int64) (result []*entity.TodoList, err error)
	GetByID(ctx context.Context, ID int64) (result *entity.TodoList, err error)
	Create(ctx context.Context, dbTrx TrxObj, params *entity.TodoList, nonZeroVal bool) error
	LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.TodoList, err error)
	Update(ctx context.Context, dbTrx TrxObj, params *entity.TodoList, changes *entity.TodoList) (err error)
	DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error
}

type TodoListRepository struct {
	GormTrxSupport
}

func NewTodoListRepository(mysql *config.Mysql) *TodoListRepository {
	return &TodoListRepository{GormTrxSupport{db: mysql.DB}}
}

func (r *TodoListRepository) GetByUserID(ctx context.Context, userID int64) (result []*entity.TodoList, err error) {
	funcName := "TodoListRepository.GetByUserID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM todo_lists WHERE user_id = ?", userID).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *TodoListRepository) GetByID(ctx context.Context, ID int64) (result *entity.TodoList, err error) {
	funcName := "TodoListRepository.GetByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.db.Raw("SELECT * FROM todo_lists WHERE id = ? LIMIT 1", ID).Scan(&result).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *TodoListRepository) Create(ctx context.Context, dbTrx TrxObj, params *entity.TodoList, nonZeroVal bool) error {
	funcName := "TodoListRepository.Create"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	cols := helper.NonZeroCols(params, nonZeroVal)
	return r.Trx(dbTrx).Select(cols).Create(&params).Error
}

func (r *TodoListRepository) LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (result *entity.TodoList, err error) {
	funcName := "TodoListRepository.LockByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	err = r.Trx(dbTrx).
		Raw("SELECT * FROM todo_lists WHERE id = ? FOR UPDATE", ID).
		Scan(&result).Error

	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrRecordNotFound()
	}

	return result, err
}

func (r *TodoListRepository) Update(ctx context.Context, dbTrx TrxObj, params *entity.TodoList, changes *entity.TodoList) (err error) {
	funcName := "TodoListRepository.Update"

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

func (r *TodoListRepository) DeleteByID(ctx context.Context, dbTrx TrxObj, id int64) error {
	funcName := "TodoListRepository.DeleteByID"

	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	err := r.Trx(dbTrx).Where("id = ?", id).Delete(&entity.TodoList{}).Error
	if err != nil {
		return err
	}

	return nil
}
