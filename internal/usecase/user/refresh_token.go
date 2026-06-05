package user_usecase

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
)

func (w *UserUsecase) RefreshToken(ctx context.Context, oldToken string) (*entity.RefreshTokenResponse, error) {
	funcName := "UserUsecase.RefreshToken"

	newToken, err := w.jwtAuth.RefreshToken(oldToken)
	if err != nil {
		helper.Log(entity.LogError, "jwtAuth.RefreshToken", funcName, err, entity.CaptureFields{}, "")

		return nil, err
	}

	return &entity.RefreshTokenResponse{Token: newToken}, nil
}
