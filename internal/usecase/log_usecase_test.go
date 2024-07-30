package usecase_test

import (
	"fmt"
	"testing"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LogUsecaseTestSuite struct {
	suite.Suite

	usecase   usecase.LogUsecase
	queue     *mocks.Queue
	zapLogger *zap.Logger
}

func (s *LogUsecaseTestSuite) SetupTest() {
	s.queue = &mocks.Queue{}
	s.zapLogger, _ = config.NewZapLog("dev")

	s.usecase = usecase.NewLogUsecase(s.queue, s.zapLogger)
}

func TestLogUsecase(t *testing.T) {
	suite.Run(t, new(LogUsecaseTestSuite))
}

func (s *LogUsecaseTestSuite) TestLog() {
	captureFieldError := map[string]string{"test": "test"}

	testcases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.queue.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name: "error publish queue",
			mockFunc: func() {
				s.queue.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			s.usecase.Log(entity.LogError, tt.name, "test.Test", fmt.Errorf("TEST"), captureFieldError, "")
		})
	}
}
