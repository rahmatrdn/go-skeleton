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

type WalletHandlerTestSuite struct {
	suite.Suite
	walletUsecase *mocks.WalletUsecase
	presenter     *mocks.Presenter
	parser        *mocks.Parser
	handler       *handler.WalletHandler
}

func (s *WalletHandlerTestSuite) SetupTest() {
	s.walletUsecase = &mocks.WalletUsecase{}
	s.presenter = &mocks.Presenter{}
	s.parser = &mocks.Parser{}

	s.handler = handler.NewWalletHandler(s.parser, s.presenter, s.walletUsecase)
}

func TestWalletHandler(t *testing.T) {
	suite.Run(t, new(WalletHandlerTestSuite))
}

func (s *WalletHandlerTestSuite) TestRegister() {
	app := fiber.New()

	s.handler.Register(app)
}

func (s *WalletHandlerTestSuite) TestGetByID() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	ID := int64(1)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail get id from parser param",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail wallet usecase GetByID",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("GetByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.GetByID(c)

			if err != nil {
				t.Errorf("GetByID() error = %v", err)
				return
			}
		})
	}
}

func (s *WalletHandlerTestSuite) TestGetByUserID() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	ID := int64(1)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.parser.On("ParserUserID", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("GetByUserID", mock.Anything, mock.Anything).Return(nil, nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail get id from parser param",
			mockFunc: func() {
				s.parser.On("ParserUserID", mock.Anything).Return(ID, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail wallet usecase GetByUserID",
			mockFunc: func() {
				s.parser.On("ParserUserID", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("GetByUserID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.GetByUserID(c)

			if err != nil {
				t.Errorf("GetByUserID() error = %v", err)
				return
			}
		})
	}
}

func (s *WalletHandlerTestSuite) TestCreate() {
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
				s.parser.On("ParserBodyRequestWithUserID", mock.Anything, mock.Anything).Return(nil).Once()
				s.walletUsecase.On("Create", mock.Anything, mock.Anything).Return(nil, nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail ParserBodyRequestWithUserID",
			mockFunc: func() {
				s.parser.On("ParserBodyRequestWithUserID", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail wallet usecase Create",
			mockFunc: func() {
				s.parser.On("ParserBodyRequestWithUserID", mock.Anything, mock.Anything).Return(nil).Once()
				s.walletUsecase.On("Create", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.Create(c)

			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}
		})
	}
}

func (s *WalletHandlerTestSuite) TestUpdate() {
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
				s.parser.On("ParserBodyWithIntIDPathParamsAndUserID", mock.Anything, mock.Anything).Return(nil).Once()
				s.walletUsecase.On("UpdateByID", mock.Anything, mock.Anything).Return(nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail wallet usecase UpdateByID",
			mockFunc: func() {
				s.parser.On("ParserBodyWithIntIDPathParamsAndUserID", mock.Anything, mock.Anything).Return(nil).Once()
				s.walletUsecase.On("UpdateByID", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail ParserBodyWithIntIDPathParamsAndUserID",
			mockFunc: func() {
				s.parser.On("ParserBodyWithIntIDPathParamsAndUserID", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.Update(c)

			if err != nil {
				t.Errorf("Update() error = %v", err)
				return
			}
		})
	}
}

func (s *WalletHandlerTestSuite) TestDelete() {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})

	defer app.ReleaseCtx(c)

	ID := int64(1)

	testCases := []struct {
		name     string
		mockFunc func()
	}{
		{
			name: "success",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("DeleteByID", mock.Anything, mock.Anything).Return(nil).Once()
				s.presenter.On("BuildSuccess", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail usecase",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, nil).Once()
				s.walletUsecase.On("DeleteByID", mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "fail parser",
			mockFunc: func() {
				s.parser.On("ParserIntIDFromPathParams", mock.Anything).Return(ID, fmt.Errorf("ERROR")).Once()
				s.presenter.On("BuildError", mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.handler.Delete(c)

			if err != nil {
				t.Errorf("Delete() error = %v", err)
				return
			}
		})
	}
}
