package todo_list_usecase

import (
	"context"

	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
)

type TodoListUsecase struct {
	todoListRepo mysql.ITodoListRepository
}

func NewTodoListUsecase(
	todoListRepo mysql.ITodoListRepository,
) *TodoListUsecase {
	return &TodoListUsecase{todoListRepo}
}

type ITodoListUsecase interface {
	GetByUserID(ctx context.Context, userID int64) (res []*entity.TodoListResponse, err error)
	GetByID(ctx context.Context, todoListID int64) (*entity.TodoListResponse, error)
	Create(ctx context.Context, todoListReq entity.TodoListReq) (*entity.TodoListResponse, error)
	UpdateByID(ctx context.Context, todoListReq entity.TodoListReq) error
	DeleteByID(ctx context.Context, todoListID int64) error
}
