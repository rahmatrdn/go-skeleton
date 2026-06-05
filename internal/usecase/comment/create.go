package comment_usecase

import (
	"context"
	"fmt"
	"time"

	errwrap "github.com/pkg/errors"
	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"
)

func (t *CommentUsecase) Create(ctx context.Context, req entity.CommentReq) (*entity.CommentResponse, error) {
	funcName := "CommentUsecase.Create"
	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(req.UserID),
		"payload": helper.ToString(req),
	}

	if errMsg := usecase.ValidateStruct(req); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(generalEntity.INVALID_PAYLOAD_CODE), errMsg)
	}

	// Validate that the post exists
	post, err := t.postRepo.GetByID(ctx, req.PostID)
	if err != nil || post == nil {
		return nil, fmt.Errorf("POST_NOT_FOUND")
	}

	userID := req.UserID
	commentPayload := &mentity.Comment{
		PostID:    req.PostID,
		UserID:    &userID,
		ParentID:  req.ParentID,
		Body:      req.Body,
		CreatedAt: time.Now(),
	}

	err = t.commentRepo.Create(ctx, nil, commentPayload, false)
	if err != nil {
		helper.LogError("commentRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.CommentResponse{
		ID:        commentPayload.ID,
		PostID:    commentPayload.PostID,
		UserID:    commentPayload.UserID,
		ParentID:  commentPayload.ParentID,
		Body:      commentPayload.Body,
		CreatedAt: commentPayload.CreatedAt,
	}, nil
}
