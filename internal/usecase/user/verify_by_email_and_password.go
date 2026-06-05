package user_usecase

import (
	"context"
	"fmt"

	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
)

func (w *UserUsecase) VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error) {
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
