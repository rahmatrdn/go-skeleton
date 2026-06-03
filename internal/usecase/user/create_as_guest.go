package user_usecase

import (
	"context"
	"fmt"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

func (w *UserUsecase) CreateAsGuest(ctx context.Context, createUserReq *entity.CreateUserReq) (*entity.CreateUserResponse, error) {
	funcName := "UserUsecase.Create"
	captureFieldError := entity.CaptureFields{
		"name": createUserReq.Name,
	}

	if errMsg := usecase.ValidateStruct(*createUserReq); errMsg != "" {
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
