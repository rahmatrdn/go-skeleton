package usecase

import (
	"context"
	"fmt"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo mysql.UserRepository
	jwtAuth  auth.JWTAuth
}

func NewUserUsecase(
	userRepo mysql.UserRepository,
	jwtAuth auth.JWTAuth,
) *User {
	return &User{userRepo, jwtAuth}
}

type UserUsecase interface {
	VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error)
	CreateAsGuest(ctx context.Context, createUserReq *entity.CreateUserReq) (*entity.CreateUserResponse, error)
}

func (w *User) VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error) {
	funcName := "UserUsecase.VerifyByEmailAndPassword"
	captureFieldError := map[string]string{"email": fmt.Sprint(req.Email)}

	user, err := w.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		helper.Log(entity.LogError, "userRepo.GetByEmail", funcName, err, captureFieldError, "")

		if err == apperr.ErrUserNotFound() {
			return nil, apperr.ErrInvalidEmailOrPassword()
		}

		return nil, err
	}

	if !helper.VerifyBcryptHash(req.Password, user.Password) {
		return nil, apperr.ErrInvalidEmailOrPassword()
	}

	token, err := w.jwtAuth.GenerateToken(user)
	if err != nil {
		helper.Log(entity.LogError, "userRepo.GenerateToken", funcName, err, captureFieldError, "")

		return nil, err
	}

	loginRes = &entity.LoginResponse{
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		RoleAccess: user.Role,
		Token:      token,
	}

	return loginRes, nil
}

func (w *User) CreateAsGuest(ctx context.Context, createUserReq *entity.CreateUserReq) (*entity.CreateUserResponse, error) {
	funcName := "UserUsecase.Create"
	captureFieldError := entity.CaptureFields{
		"name": createUserReq.Name,
	}

	if errMsg := ValidateStruct(*createUserReq); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(entity.INVALID_PAYLOAD_CODE), errMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.LogError("bcrypt.GenerateFromPassword", funcName, err, captureFieldError, "")

		return nil, err
	}

	user := &mentity.User{
		Name:     createUserReq.Name,
		Email:    createUserReq.Email,
		Role:     int8(entity.Guest),
		Phone:    createUserReq.Phone,
		Password: string(hashedPassword),
	}

	err = w.userRepo.Create(ctx, nil, user)
	if err != nil {
		helper.LogError("userRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	token, err := w.jwtAuth.GenerateToken(user)
	if err != nil {
		helper.LogError("userRepo.GetByEmail", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.CreateUserResponse{
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		RoleAccess: entity.GetRoleName(entity.UserRole(user.Role)),
		Token:      token,
	}, nil
}
