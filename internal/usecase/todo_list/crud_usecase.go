package todo_list_usecase

import (
	"context"
	"fmt"
	"time"

	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/entity"
	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
)

type CrudUsecase struct {
	todoListRepo mysql.TodoListRepository
}

func NewCrudUsecase(
	todoListRepo mysql.TodoListRepository,
) *CrudUsecase {
	return &CrudUsecase{todoListRepo}
}

type ICrudUsecase interface {
	GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error)
	GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error)
	Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error)
	UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error
	DeleteByID(ctx context.Context, todoListID int64) error
}

func (t *CrudUsecase) GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error) {
	funcName := "CrudUsecase.GetByUserID"
	captureFieldError := entity.CaptureFields{
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

func (t *CrudUsecase) GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error) {
	funcName := "CrudUsecase.GetByID"
	captureFieldError := entity.CaptureFields{
		"user_id": helper.ToString(todoListID),
	}

	data, err := t.todoListRepo.GetByID(ctx, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.GetByID", funcName, err, captureFieldError, "")

		return nil, err
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

func (t *CrudUsecase) Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error) {
	funcName := "CrudUsecase.Create"
	captureFieldError := entity.CaptureFields{
		"user_id": helper.ToString(todoListReq.UserID),
		"payload": helper.ToString(todoListReq),
	}

	if errMsg := usecase.ValidateStruct(todoListReq); errMsg != "" {
		return nil, errwrap.Wrap(fmt.Errorf(entity.INVALID_PAYLOAD_CODE), errMsg)
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

func (t *CrudUsecase) UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error {
	funcName := "CrudUsecase.UpdateByID"
	todoListID := todoListReq.ID

	captureFieldError := entity.CaptureFields{
		"user_id": helper.ToString(todoListReq.UserID),
		"payload": helper.ToString(todoListReq),
	}

	if err := mysql.DBTransaction(t.todoListRepo, func(trx mysql.TrxObj) error {
		lockedData, err := t.todoListRepo.LockByID(ctx, trx, todoListID)
		if err != nil {
			helper.LogError("todoListRepo.LockByID", funcName, err, captureFieldError, "")

			return err
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

func (t *CrudUsecase) DeleteByID(ctx context.Context, todoListID int64) error {
	funcName := "CrudUsecase.DeleteByID"
	captureFieldError := entity.CaptureFields{
		"todo_list_id": helper.ToString(todoListID),
	}

	err := t.todoListRepo.DeleteByID(ctx, nil, todoListID)
	if err != nil {
		helper.LogError("todoListRepo.DeleteByID", funcName, err, captureFieldError, "")

		return err
	}

	return nil
}
