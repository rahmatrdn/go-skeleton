package comment_usecase

import (
	"context"
	"fmt"
	"time"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/comment/entity"
)

func (t *CommentUsecase) UpdateByID(ctx context.Context, req entity.CommentReq) error {
	funcName := "CommentUsecase.UpdateByID"
	commentID := req.ID

	captureFieldError := generalEntity.CaptureFields{
		"comment_id": helper.ToString(commentID),
		"payload":    helper.ToString(req),
	}

	if err := mysql.DBTransaction(t.commentRepo, func(trx mysql.TrxObj) error {
		lockedData, err := t.commentRepo.LockByID(ctx, trx, commentID)
		if err != nil {
			helper.LogError("commentRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}
		if lockedData == nil {
			return fmt.Errorf("DATA IS NOT EXIST")
		}

		if err := t.commentRepo.Update(ctx, trx, lockedData, &mentity.Comment{
			Body:      req.Body,
			UpdatedAt: time.Now(),
		}); err != nil {
			helper.LogError("commentRepo.Update", funcName, err, captureFieldError, "")

			return err
		}

		return nil
	}); err != nil {
		helper.LogError("commentRepo.DBTransaction", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
