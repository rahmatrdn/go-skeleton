package comment_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"
)

func (t *CommentUsecase) GetByPostID(ctx context.Context, postID int64) ([]*entity.CommentResponse, error) {
	funcName := "CommentUsecase.GetByPostID"
	captureFieldError := generalEntity.CaptureFields{
		"post_id": helper.ToString(postID),
	}

	data, err := t.commentRepo.GetByPostID(ctx, postID)
	if err != nil {
		helper.LogError("commentRepo.GetByPostID", funcName, err, captureFieldError, "")

		return nil, err
	}

	var result []*entity.CommentResponse
	for _, d := range data {
		result = append(result, &entity.CommentResponse{
			ID:        d.ID,
			PostID:    d.PostID,
			UserID:    d.UserID,
			ParentID:  d.ParentID,
			Body:      d.Body,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		})
	}

	return result, nil
}
