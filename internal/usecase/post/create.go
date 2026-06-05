package post_usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	errwrap "github.com/pkg/errors"
	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/post/entity"
)

func slugify(title string) string {
	s := strings.ToLower(title)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func generateUniqueSlug(title string) string {
	base := slugify(title)
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	return fmt.Sprintf("%s-%s", base, suffix)
}

func (t *PostUsecase) Create(ctx context.Context, req entity.PostReq) (*entity.PostResponse, error) {
	funcName := "PostUsecase.Create"
	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(req.UserID),
		"payload": helper.ToString(req),
	}

	if errMsg := usecase.ValidateStruct(req); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(generalEntity.INVALID_PAYLOAD_CODE), errMsg)
	}

	now := time.Now()
	var publishedAt *time.Time
	if req.Status == entity.PostStatusPublished {
		publishedAt = &now
	}

	postPayload := &mentity.Post{
		UserID:      req.UserID,
		Title:       req.Title,
		Slug:        generateUniqueSlug(req.Title),
		Body:        req.Body,
		Status:      req.Status,
		PublishedAt: publishedAt,
		CreatedAt:   now,
	}

	err := t.postRepo.Create(ctx, nil, postPayload, false)
	if err != nil {
		helper.LogError("postRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.PostResponse{
		ID:          postPayload.ID,
		UserID:      postPayload.UserID,
		Title:       postPayload.Title,
		Slug:        postPayload.Slug,
		Body:        postPayload.Body,
		Status:      postPayload.Status,
		PublishedAt: postPayload.PublishedAt,
		CreatedAt:   postPayload.CreatedAt,
	}, nil
}
