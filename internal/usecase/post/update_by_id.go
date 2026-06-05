package post_usecase

import (
	"context"
	"fmt"
	"time"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"
)

func (t *PostUsecase) UpdateByID(ctx context.Context, req entity.PostReq) error {
	funcName := "PostUsecase.UpdateByID"
	postID := req.ID

	captureFieldError := generalEntity.CaptureFields{
		"post_id": helper.ToString(postID),
		"payload": helper.ToString(req),
	}

	if err := mysql.DBTransaction(t.postRepo, func(trx mysql.TrxObj) error {
		lockedData, err := t.postRepo.LockByID(ctx, trx, postID)
		if err != nil {
			helper.LogError("postRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}
		if lockedData == nil {
			return fmt.Errorf("DATA IS NOT EXIST")
		}

		now := time.Now()

		// Regenerate slug if title changed
		newSlug := lockedData.Slug
		if req.Title != lockedData.Title {
			newSlug = generateUniqueSlug(req.Title)
		}

		// Set publishedAt when status changes to Published
		var publishedAt *time.Time
		if req.Status == entity.PostStatusPublished && lockedData.PublishedAt == nil {
			publishedAt = &now
		} else {
			publishedAt = lockedData.PublishedAt
		}

		if err := t.postRepo.Update(ctx, trx, lockedData, &mentity.Post{
			Title:       req.Title,
			Slug:        newSlug,
			Body:        req.Body,
			Status:      req.Status,
			PublishedAt: publishedAt,
			UpdatedAt:   now,
		}); err != nil {
			helper.LogError("postRepo.Update", funcName, err, captureFieldError, "")

			return err
		}

		return nil
	}); err != nil {
		helper.LogError("postRepo.DBTransaction", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
