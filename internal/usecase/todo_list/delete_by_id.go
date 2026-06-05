package todo_list_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
)

func (t *TodoListUsecase) DeleteByID(ctx context.Context, todoListID int64) error {
	funcName := "TodoListUsecase.DeleteByID"
	captureFieldError := generalEntity.CaptureFields{
		"todo_list_id": helper.ToString(todoListID),
	}

	err := t.todoListRepo.DeleteByID(ctx, nil, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
