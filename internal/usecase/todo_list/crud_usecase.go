package todo_list_usecase

import (
	"context"
	"fmt"
	"time"

	errwrap "github.com/pkg/errors"
	generalEntity "github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
)

type CrudTodoListUsecase struct {
	todoListRepo mysql.ITodoListRepository
}

func NewCrudTodoListUsecase(
	todoListRepo mysql.ITodoListRepository,
) *CrudTodoListUsecase {
	return &CrudTodoListUsecase{todoListRepo}
}

type ICrudTodoListUsecase interface {
	GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error)
	GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error)
	Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error)
	UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error
	DeleteByID(ctx context.Context, todoListID int64) error
}

func (t *CrudTodoListUsecase) GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error) {
	funcName := "CrudTodoListUsecase.GetByUserID"
	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(userID),
	}

	result, err := t.todoListRepo.GetByUserID(ctx, userID)
	if err != nil {
		helper.LogError("todoListRepo.GetByUserID", funcName, err, captureFieldError, "")

		return nil, err
	}

	for _, v := range result {
		res = append(res, &entity.TodoListResponse{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			DoingAt:     helper.ConvertToJakartaDate(v.DoingAt),
			CreatedAt:   helper.ConvertToJakartaTime(v.CreatedAt),
			UpdatedAt:   helper.ConvertToJakartaTime(v.UpdatedAt),
		})
	}

	return res, nil
}

func (t *CrudTodoListUsecase) GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error) {
	funcName := "CrudTodoListUsecase.GetByID"
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

func (t *CrudTodoListUsecase) Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error) {
	funcName := "CrudTodoListUsecase.Create"
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
		DoingAt:     helper.ConvertToJakartaDate(todoListPayload.DoingAt),
		CreatedAt:   helper.ConvertToJakartaTime(todoListPayload.CreatedAt),
	}, nil
}

func (t *CrudTodoListUsecase) UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error {
	funcName := "CrudTodoListUsecase.UpdateByID"
	todoListID := todoListReq.ID

	captureFieldError := generalEntity.CaptureFields{
		"user_id": helper.ToString(todoListReq.UserID),
		"payload": helper.ToString(todoListReq),
	}

	if err := mysql.DBTransaction(t.todoListRepo, func(trx mysql.TrxObj) error {
		lockedData, err := t.todoListRepo.LockByID(ctx, trx, todoListID)
		if err != nil {
			helper.LogError("todoListRepo.LockByID", funcName, err, captureFieldError, "")

			return err
		}
		if lockedData == nil {
			return fmt.Errorf("DATA IS NOT EXIST")
		}

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

func (t *CrudTodoListUsecase) DeleteByID(ctx context.Context, todoListID int64) error {
	funcName := "CrudTodoListUsecase.DeleteByID"
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
