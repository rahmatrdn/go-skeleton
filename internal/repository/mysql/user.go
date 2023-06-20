package mysql

import (
	"context"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/config"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"gorm.io/gorm"

	"gorm.io/gorm/clause"
)

type UserRepository interface {
	TrxSupportRepo
	Create(ctx context.Context, dbTrx TrxObj, user *entity.User) error
	LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByEmailAndRole(ctx context.Context, email string, role entity.RoleType) (*entity.User, error)
}

type User struct {
	GormTrxSupport
}

func NewUserRepository(mysql *config.Mysql) *User {
	return &User{GormTrxSupport{db: mysql.DB}}
}

func (u *User) Create(ctx context.Context, dbTrx TrxObj, user *entity.User) error {
	funcName := "UserRepository.Create"
	if err := helper.CheckDeadline(ctx); err != nil {
		return errwrap.Wrap(err, funcName)
	}

	return u.Trx(dbTrx).Create(&user).Error
}

func (u *User) LockByID(ctx context.Context, dbTrx TrxObj, ID int64) (*entity.User, error) {
	funcName := "UserRepository.GetUserByID"
	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	var user *entity.User
	err := u.Trx(dbTrx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ID).Take(&user).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrUserNotFound()
	}

	return user, err
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	funcName := "UserRepository.GetByEmail"
	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	var user *entity.User
	err := u.db.Where("email = ?", email).Take(&user).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrUserNotFound()
	}

	return user, err
}

func (u *User) GetByEmailAndRole(ctx context.Context, email string, role entity.RoleType) (*entity.User, error) {
	funcName := "UserRepository.GetByEmailAndRole"
	if err := helper.CheckDeadline(ctx); err != nil {
		return nil, errwrap.Wrap(err, funcName)
	}

	var user *entity.User
	err := u.db.Where("email = ? AND role = ?", email, role).Take(&user).Error
	if errwrap.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrUserNotFound()
	}

	return user, err
}
