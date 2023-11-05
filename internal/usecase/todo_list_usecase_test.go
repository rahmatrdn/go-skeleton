package usecase_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/rahmatrdn/go-skeleton/internal/helper"
	"github.com/rahmatrdn/go-skeleton/internal/usecase"
	"github.com/rahmatrdn/go-skeleton/tests/fixture/factory"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"

	"github.com/rahmatrdn/go-skeleton/entity"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TodoListUsecaseTestSuite struct {
	suite.Suite

	todoListRepo *mocks.TodoListRepository
	usecase      usecase.TodoListUsecase
	trxObj       *mocks.TrxObj

	ctx context.Context
}

func (s *TodoListUsecaseTestSuite) SetupTest() {
	s.todoListRepo = &mocks.TodoListRepository{}

	s.usecase = usecase.NewTodoListUsecase(s.todoListRepo)
	s.trxObj = &mocks.TrxObj{}

	s.todoListRepo.On("Begin").Return(s.trxObj, nil)
	s.trxObj.On("Rollback", mock.Anything).Return(nil)
	s.trxObj.On("Commit", mock.Anything).Return(nil)

	s.ctx = context.Background()
}

func TestTodoListUsecase(t *testing.T) {
	suite.Run(t, new(TodoListUsecaseTestSuite))
}

func (s *TodoListUsecaseTestSuite) TestGetByUserID() {
	todoList := factory.StubbedTodoLists()
	userID := todoList[0].UserID

	var res []*entity.TodoListResponse
	for _, v := range todoList {
		res = append(res, &entity.TodoListResponse{
			ID:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			DoingAt:     v.DoingAt,
			CreatedAt:   helper.ConvertToJakartaTime(v.CreatedAt),
			UpdatedAt:   helper.ConvertToJakartaTime(v.UpdatedAt),
		})
	}

	testcases := []struct {
		name     string
		mockFunc func()
		want     []*entity.TodoListResponse
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.todoListRepo.On("GetByUserID", s.ctx, userID).Return(todoList, nil).Once()
			},
			want: res,
		},
		{
			name: "fail get by user id",
			mockFunc: func() {
				s.todoListRepo.On("GetByUserID", s.ctx, userID).Return(nil, fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			got, err := s.usecase.GetByUserID(s.ctx, userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *TodoListUsecaseTestSuite) TestGetByID() {
	todoList := factory.StubbedTodoList()
	id := todoList.ID

	expectedResponse := &entity.TodoListResponse{
		ID:          todoList.ID,
		Title:       todoList.Title,
		Description: todoList.Description,
		DoingAt:     todoList.DoingAt,
		CreatedAt:   helper.ConvertToJakartaTime(todoList.CreatedAt),
		UpdatedAt:   helper.ConvertToJakartaTime(todoList.UpdatedAt),
	}

	testcases := []struct {
		name     string
		mockFunc func()
		want     *entity.TodoListResponse
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.todoListRepo.On("GetByID", s.ctx, id).Return(todoList, nil).Once()
			},
			want: expectedResponse,
		},
		{
			name: "fail get by user id",
			mockFunc: func() {
				s.todoListRepo.On("GetByID", s.ctx, id).Return(nil, fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			got, err := s.usecase.GetByID(s.ctx, id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *TodoListUsecaseTestSuite) TestCreate() {
	todoList := factory.StubbedTodoList()

	params := &entity.TodoListReq{
		UserID:      todoList.UserID,
		Title:       todoList.Title,
		Description: todoList.Description,
		DoingAt:     todoList.DoingAt,
	}

	testCases := []struct {
		name     string
		params   *entity.TodoListReq
		mockFunc func(params *entity.TodoListReq)
		wantErr  bool
	}{
		{
			name:   "success",
			params: params,
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("Create", s.ctx, mock.Anything, mock.Anything, false).Return(nil).Once()
			},
		},
		{
			name:   "fail on create",
			params: params,
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("Create", s.ctx, mock.Anything, mock.Anything, false).Return(fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
		{
			name: "fail on validate struct",
			params: &entity.TodoListReq{
				UserID:      todoList.UserID,
				Title:       todoList.Title,
				Description: todoList.Description,
				DoingAt:     "",
			},
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("Create", s.ctx, mock.Anything, mock.Anything, false).Return(nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.params)

			_, err := s.usecase.Create(s.ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (s *TodoListUsecaseTestSuite) TestUpdateByID() {
	todoList := factory.StubbedTodoList()
	ID := todoList.ID

	params := &entity.TodoListReq{
		UserID:      todoList.UserID,
		Title:       todoList.Title,
		Description: todoList.Description,
		DoingAt:     todoList.DoingAt,
		ID:          ID,
	}

	testCases := []struct {
		name     string
		params   *entity.TodoListReq
		mockFunc func(params *entity.TodoListReq)
		wantErr  bool
	}{
		{
			name:   "success",
			params: params,
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("LockByID", s.ctx, mock.Anything, ID).Return(todoList, nil).Once()
				s.todoListRepo.On("Update", s.ctx, mock.Anything, todoList, mock.Anything).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "fail update to db",
			params: params,
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("LockByID", s.ctx, mock.Anything, ID).Return(todoList, nil).Once()
				s.todoListRepo.On("Update", s.ctx, mock.Anything, todoList, mock.Anything).Return(fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
		{
			name:   "fail lock todo list by id",
			params: params,
			mockFunc: func(params *entity.TodoListReq) {
				s.todoListRepo.On("LockByID", s.ctx, mock.Anything, ID).Return(todoList, fmt.Errorf("ERROR")).Once()
				s.todoListRepo.On("Update", s.ctx, mock.Anything, todoList, mock.Anything).Return(nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.params)

			err := s.usecase.UpdateByID(s.ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (s *TodoListUsecaseTestSuite) TestDeleteByID() {
	todoList := factory.StubbedTodoList()

	testCases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() {
				s.todoListRepo.On("DeleteByID", s.ctx, nil, todoList.ID).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "fail",
			mockFunc: func() {
				s.todoListRepo.On("DeleteByID", s.ctx, nil, todoList.ID).Return(fmt.Errorf("ERROR")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := s.usecase.DeleteByID(s.ctx, todoList.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAddressByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
