package post_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
)

func (t *PostUsecase) DeleteByID(ctx context.Context, postID int64) error {
	funcName := "PostUsecase.DeleteByID"
	captureFieldError := generalEntity.CaptureFields{
		"post_id": helper.ToString(postID),
	}

	err := t.postRepo.DeleteByID(ctx, nil, postID)
	if err != nil {
		helper.LogError("postRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
