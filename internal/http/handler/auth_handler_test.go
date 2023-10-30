package handler_test

import (
	"fmt"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rahmatrdn/go-skeleton/internal/http/handler"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecase
	presenter   *mocks.Presenter
	parser      *mocks.Parser
	handler     *handler.AuthHandler
}

func (s *AuthHandlerTestSuite) SetupTest() {
	s.userUsecase = &mocks.UserUsecase{}
	s.presenter = &mocks.Presenter{}
	s.parser = &mocks.Parser{}

	s.handler = handler.NewAuthHandler(s.parser, s.presenter, s.userUsecase)
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (s *AuthHandlerTestSuite) TestRegister() {
	app := fiber.New()

	s.handler.Register(app)
}

func (s *AuthHandlerTestSuite) TestCreateAsGuest() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(nil).Once()
				s.userUsecase.On("CreateAsGuest", mock.Anything, mock.Anything).Return(nil, nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail get id from parser param",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail usecase",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(nil).Once()
				s.userUsecase.On("CreateAsGuest", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.CreateAsGuest(c)

			if err != nil {
				t.Errorf("Login() error = %v", err)
				return
			}
		})
	}
}

func (s *AuthHandlerTestSuite) TestLogin() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(nil).Once()
				s.userUsecase.On("VerifyByEmailAndPassword", mock.Anything, mock.Anything).Return(nil, nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail get id from parser param",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail usecase",
			mockFunc: func() {
				s.parser.On("ParserBodyRequest", mock.Anything, mock.Anything).Return(nil).Once()
				s.userUsecase.On("VerifyByEmailAndPassword", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.Login(c)

			if err != nil {
				t.Errorf("Login() error = %v", err)
				return
			}
		})
	}
}

func (s *AuthHandlerTestSuite) TestCheckToken() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.CheckToken(c)

			if err != nil {
				t.Errorf("CheckToken() error = %v", err)
				return
			}
		})
	}
}
