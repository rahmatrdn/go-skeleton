package post_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"
)

func (t *PostUsecase) GetAll(ctx context.Context) ([]*entity.PostResponse, error) {
	funcName := "PostUsecase.GetAll"

	data, err := t.postRepo.GetAll(ctx)
	if err != nil {
		helper.LogError("postRepo.GetAll", funcName, err, generalEntity.CaptureFields{}, "")

		return nil, err
	}

	var result []*entity.PostResponse
	for _, d := range data {
		result = append(result, &entity.PostResponse{
			ID:          d.ID,
			UserID:      d.UserID,
			UserName:    d.UserName,
			Title:       d.Title,
			Slug:        d.Slug,
			Body:        d.Body,
			Status:      d.Status,
			PublishedAt: d.PublishedAt,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return result, nil
}
