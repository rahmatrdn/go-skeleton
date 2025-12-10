package usecase_test

import (
	"encoding/json"
	"testing"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type ValidationTestStruct struct {
	Name  string `validate:"required" name:"Nama"`
	Email string `validate:"required,email" name:"Email"`
}

func TestValidateStructProcess(t *testing.T) {
	testcases := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "Success",
			data: ValidationTestStruct{
				Name:  "Test",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name:    "Error Required",
			data:    ValidationTestStruct{},
			wantErr: true,
		},
		{
			name: "Error Email Format",
			data: ValidationTestStruct{
				Name:  "Test",
				Email: "invalid-email",
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			errs := usecase.ValidateStructProcess(tt.data)
			if tt.wantErr {
				assert.NotEmpty(t, errs)
				// Optional: Verify strict fields
				if len(errs) > 0 {
					assert.IsType(t, entity.ErrorResponse{}, errs[0])
				}
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	testcases := []struct {
		name       string
		data       interface{}
		wantSuffix string
		wantEmpty  bool
	}{
		{
			name:       "Success",
			data:       ValidationTestStruct{Name: "Test", Email: "test@example.com"},
			wantSuffix: "",
			wantEmpty:  true,
		},
		{
			name:       "Error",
			data:       ValidationTestStruct{},
			wantSuffix: "XX",
			wantEmpty:  false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			msg := usecase.ValidateStruct(tt.data)
			if tt.wantEmpty {
				assert.Empty(t, msg)
			} else {
				assert.NotEmpty(t, msg)
				assert.Contains(t, msg, tt.wantSuffix)

				// Verify it's a valid JSON + XX
				jsonPart := msg[:len(msg)-2]
				var errs []entity.ErrorResponse
				err := json.Unmarshal([]byte(jsonPart), &errs)
				assert.NoError(t, err)
				assert.NotEmpty(t, errs)
			}
		})
	}
}
