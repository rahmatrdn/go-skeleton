package todo_list_usecase

import (
	"context"
	"fmt"
	"time"

	errwrap "github.com/pkg/errors"
	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
)

func (t *TodoListUsecase) Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error) {
	funcName := "TodoListUsecase.Create"
	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(todoListReq.UserID),
		"payload": helper.ToString(todoListReq),
	}

	if errMsg := usecase.ValidateStruct(todoListReq); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(generalEntity.INVALID_PAYLOAD_CODE), errMsg)
	}

	doingAt, _ := helper.ParseDate(todoListReq.DoingAt)

	todoListPayload := &mentity.TodoList{
		UserID:      todoListReq.UserID,
		Title:       todoListReq.Title,
		Description: todoListReq.Description,
		DoingAt:     doingAt,
		CreatedAt:   time.Now(),
	}

	err := t.todoListRepo.Create(ctx, nil, todoListPayload, false)
	if err != nil {
		helper.LogError("todoListRepo.Create", funcName, err, captureFieldError, "")

		return nil, err
	}

	return &entity.TodoListResponse{
		ID:          todoListPayload.ID,
		Title:       todoListPayload.Title,
		Description: todoListPayload.Description,
		DoingAt:     todoListPayload.DoingAt,
		CreatedAt:   todoListPayload.CreatedAt,
	}, nil
}
