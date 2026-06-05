package comment_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"
)

func (t *CommentUsecase) GetByID(ctx context.Context, commentID int64) (*entity.CommentResponse, error) {
	funcName := "CommentUsecase.GetByID"
	captureFieldError := generalEntity.CaptureFields{
		"comment_id": helper.ToString(commentID),
	}

	data, err := t.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		helper.LogError("commentRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &entity.CommentResponse{
		ID:        data.ID,
		PostID:    data.PostID,
		UserID:    data.UserID,
		ParentID:  data.ParentID,
		Body:      data.Body,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
}
