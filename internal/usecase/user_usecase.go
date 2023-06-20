package usecase

import (
	"context"
	"fmt"

	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo mysql.UserRepository
	logger   LogUsecase
	jwtAuth  auth.JWT
}

func NewUserUsecase(
	userRepo mysql.UserRepository,
	logger LogUsecase,
	jwtAuth auth.JWT,
) *User {
	return &User{userRepo, logger, jwtAuth}
}

type UserUsecase interface {
	VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error)
}

func (w *User) VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error) {
	funcName := "UserUsecase.VerifyByEmailAndPassword"
	captureFieldError := map[string]string{"email": fmt.Sprint(req.Email)}

	user, err := w.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.GetByEmail", funcName, err, captureFieldError, "")

		if err == apperr.ErrUserNotFound() {
			return nil, apperr.ErrInvalidEmailOrPassword()
		}

		return nil, err
	}

	if !VerifyBcryptHash(req.Password, user.Password) {
		return nil, apperr.ErrInvalidEmailOrPassword()
	}

	loginRes = &entity.LoginResponse{
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		RoleAccess: user.Role,
	}

	token, err := w.jwtAuth.GenerateToken(*loginRes)
	if err != nil {
		w.logger.Log(entity.LogError, "walletRepo.GetByEmail", funcName, err, captureFieldError, "")

		return nil, err
	}

	loginRes.Token = token

	return loginRes, nil
}

// Fungsi ini digunakan untuk memverifikasi apakah hash bcrypt cocok dengan plaintext
func VerifyBcryptHash(plaintext, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
