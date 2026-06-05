package user_usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rahmatrdn/go-skeleton/entity"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	user_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/user"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	usecase *user_usecase.UserUsecase
	repo    *mocks.UserRepository
	jwtAuth *mocks.JWTAuth
}

func (s *UserUsecaseTestSuite) SetupTest() {
	s.repo = &mocks.UserRepository{}
	s.jwtAuth = &mocks.JWTAuth{}
	s.usecase = user_usecase.NewUserUsecase(s.repo, s.jwtAuth)
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (s *UserUsecaseTestSuite) TestVerifyByEmailAndPassword() {
	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := &entity.LoginReq{
		Email:    "test@example.com",
		Password: password,
	}

	testcases := []struct {
		name         string
		mockFunc     func()
		wantToken    string
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			mockFunc: func() {
				user := &mentity.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Name:     "Test User",
					Role:     1,
				}
				s.repo.On("GetByEmail", ctx, req.Email).Return(user, nil).Once()
				s.jwtAuth.On("GenerateToken", user).Return("valid_token", nil).Once()
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name: "User Not Found",
			mockFunc: func() {
				s.repo.On("GetByEmail", ctx, req.Email).Return(nil, apperr.ErrUserNotFound()).Once()
			},
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrInvalidEmailOrPassword(), err)
			},
		},
		{
			name: "Invalid Password",
			mockFunc: func() {
				otherHash, _ := bcrypt.GenerateFromPassword([]byte("wrongpass"), bcrypt.DefaultCost)
				userWithWrongPass := &mentity.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(otherHash),
				}
				s.repo.On("GetByEmail", ctx, req.Email).Return(userWithWrongPass, nil).Once()
			},
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrInvalidEmailOrPassword(), err)
			},
		},
		{
			name: "Repo Error",
			mockFunc: func() {
				s.repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("db error")).Once()
			},
			wantErr: true,
		},
		{
			name: "Token Generation Error",
			mockFunc: func() {
				user := &mentity.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				s.repo.On("GetByEmail", ctx, req.Email).Return(user, nil).Once()
				s.jwtAuth.On("GenerateToken", user).Return("", errors.New("token error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			res, err := s.usecase.VerifyByEmailAndPassword(ctx, req)
			if tt.wantErr {
				s.Error(err)
				s.Nil(res)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				if res != nil {
					s.Equal(tt.wantToken, res.Token)
				} else {
					s.Fail("Expected response but got nil")
				}
			}
		})
	}
}

func (s *UserUsecaseTestSuite) TestCreateAsGuest() {
	ctx := context.Background()
	req := &entity.CreateUserReq{
		Name:            "Guest User",
		Email:           "guest@example.com",
		Password:        "password123",
		ReenterPassword: "password123",
		Phone:           "08123456789",
		RoleAccess:      2,
	}

	testcases := []struct {
		name     string
		req      *entity.CreateUserReq
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success",
			req:  req,
			mockFunc: func() {
				s.repo.On("Create", ctx, mock.Anything, mock.Anything).Return(nil).Once()
				s.jwtAuth.On("GenerateToken", mock.Anything).Return("guest_token", nil).Once()
			},
			wantErr: false,
		},
		{
			name:     "Validation Error",
			req:      &entity.CreateUserReq{},
			mockFunc: func() {},
			wantErr:  true,
		},
		{
			name: "Bcrypt Error",
			req: &entity.CreateUserReq{
				Name:            "Guest User",
				Email:           "guest@example.com",
				Password:        string(make([]byte, 73)),
				ReenterPassword: string(make([]byte, 73)),
				Phone:           "08123456789",
				RoleAccess:      2,
			},
			mockFunc: func() {},
			wantErr:  true,
		},
		{
			name: "Repo Create Error",
			req:  req,
			mockFunc: func() {
				s.repo.On("Create", ctx, mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			wantErr: true,
		},
		{
			name: "Token Generation Error",
			req:  req,
			mockFunc: func() {
				s.repo.On("Create", ctx, mock.Anything, mock.Anything).Return(nil).Once()
				s.jwtAuth.On("GenerateToken", mock.Anything).Return("", errors.New("token error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			res, err := s.usecase.CreateAsGuest(ctx, tt.req)
			if tt.wantErr {
				s.Error(err)
				s.Nil(res)
			} else {
				s.NoError(err)
				if res != nil {
					s.Equal("guest_token", res.Token)
				} else {
					s.Fail("Expected response but got nil")
				}
			}
		})
	}
}
