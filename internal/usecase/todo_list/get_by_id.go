package todo_list_usecase

import (
	"context"

	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
)

func (t *TodoListUsecase) GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error) {
	funcName := "TodoListUsecase.GetByID"
	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(todoListID),
	}

	data, err := t.todoListRepo.GetByID(ctx, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &entity.TodoListResponse{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		DoingAt:     helper.ConvertToJakartaDate(data.DoingAt),
		CreatedAt:   helper.ConvertToJakartaTime(data.CreatedAt),
		UpdatedAt:   helper.ConvertToJakartaTime(data.UpdatedAt),
	}, nil
}
