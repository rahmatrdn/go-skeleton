package user_usecase

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/http/auth"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
)

type UserUsecase struct {
	userRepo mysql.UserRepository
	jwtAuth  auth.JWTAuth
}

func NewUserUsecase(
	userRepo mysql.UserRepository,
	jwtAuth auth.JWTAuth,
) *UserUsecase {
	return &UserUsecase{userRepo, jwtAuth}
}

type IUserUsecase interface {
	VerifyByEmailAndPassword(ctx context.Context, req *entity.LoginReq) (loginRes *entity.LoginResponse, err error)
	CreateAsGuest(ctx context.Context, createUserReq *entity.CreateUserReq) (*entity.CreateUserResponse, error)
}
