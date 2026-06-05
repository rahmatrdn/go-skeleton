package comment_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
)

func (t *CommentUsecase) DeleteByID(ctx context.Context, commentID int64) error {
	funcName := "CommentUsecase.DeleteByID"
	captureFieldError := generalEntity.CaptureFields{
		"comment_id": helper.ToString(commentID),
	}

	err := t.commentRepo.DeleteByID(ctx, nil, commentID)
	if err != nil {
		helper.LogError("commentRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
