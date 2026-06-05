package post_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"
)

func (t *PostUsecase) GetByID(ctx context.Context, postID int64) (*entity.PostResponse, error) {
	funcName := "PostUsecase.GetByID"
	captureFieldError := generalEntity.CaptureFields{
		"post_id": helper.ToString(postID),
	}

	data, err := t.postRepo.GetByID(ctx, postID)
	if err != nil {
		helper.LogError("postRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &entity.PostResponse{
		ID:          data.ID,
		UserID:      data.UserID,
		Title:       data.Title,
		Slug:        data.Slug,
		Body:        data.Body,
		Status:      data.Status,
		PublishedAt: data.PublishedAt,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}, nil
}
