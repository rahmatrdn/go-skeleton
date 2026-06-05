package todo_list_usecase

import (
	"context"
	"fmt"
	"time"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
)

func (t *TodoListUsecase) UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error {
	funcName := "TodoListUsecase.UpdateByID"
	todoListID := todoListReq.ID

	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(todoListReq.UserID),
		"payload": helper.ToString(todoListReq),
	}

	// Start DB Transaction
	if err := mysql.DBTransaction(t.todoListRepo, func(trx mysql.TrxObj) error {
		// Locking Data
		lockedData, err := t.todoListRepo.LockByID(ctx, trx, todoListID)
		if err != nil {
			helper.LogError("todoListRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}
		if lockedData == nil {
			return fmt.Errorf("DATA IS NOT EXIST")
		}

		// Process Update
		doingAt, _ := helper.ParseDate(todoListReq.DoingAt)
		if err := t.todoListRepo.Update(ctx, trx, lockedData, &mentity.TodoList{
			Title:       todoListReq.Title,
			Description: todoListReq.Description,
			DoingAt:     doingAt,
			UpdatedAt:   time.Now(),
		}); err != nil {
			helper.LogError("todoListRepo.Update", funcName, err, captureFieldError, "")

			return err
		}

		return nil
	}); err != nil {
		helper.LogError("todoListRepo.DBTransaction", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
