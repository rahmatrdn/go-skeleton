package usecase_test

import (
	"testing"

	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/stretchr/testify/suite"
)

type ValidatorUsecaseTestSuite struct {
	suite.Suite
	usecase *usecase.Validator
}

func (s *ValidatorUsecaseTestSuite) SetupTest() {
	s.usecase = usecase.NewValidatorUsecase()
}

func TestValidatorUsecase(t *testing.T) {
	suite.Run(t, new(ValidatorUsecaseTestSuite))
}

type ValidatorTestStruct struct {
	Name  string `validate:"required" name:"Nama"`
	Email string `validate:"required,email" name:"Email"`
}

func (s *ValidatorUsecaseTestSuite) TestValidate() {
	testcases := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "Success",
			data: ValidatorTestStruct{
				Name:  "Test",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name:    "Error Required",
			data:    ValidatorTestStruct{},
			wantErr: true,
		},
		{
			name: "Error Email Format",
			data: ValidatorTestStruct{
				Name:  "Test",
				Email: "invalid-email",
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			errs := s.usecase.Validate(tt.data)
			if tt.wantErr {
				s.NotEmpty(errs)
			} else {
				s.Empty(errs)
			}
		})
	}
}

func (s *ValidatorUsecaseTestSuite) TestValidateWithMessage() {
	testcases := []struct {
		name       string
		data       interface{}
		wantSuffix string
	}{
		{
			name:       "Success",
			data:       ValidatorTestStruct{Name: "Test", Email: "test@example.com"},
			wantSuffix: "XX",
		},
		{
			name:       "Error",
			data:       ValidatorTestStruct{},
			wantSuffix: "XX",
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			msg := s.usecase.ValidateWithMessage(tt.data)
			s.Contains(msg, tt.wantSuffix)
			s.NotEmpty(msg)
		})
	}
}
