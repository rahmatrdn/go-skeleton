package todo_list_usecase_test

import (
	"context"
	"errors"
	"testing"

	mentity "github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	todo_list_usecase "github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list"
	"github.com/rahmatrdn/go-skeleton/internal/usecase/todo_list/entity"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CrudTodoListUsecaseTestSuite struct {
	suite.Suite
	usecase *todo_list_usecase.CrudTodoListUsecase
	repo    *mocks.ITodoListRepository
	trxObj  *mocks.TrxObj
}

func (s *CrudTodoListUsecaseTestSuite) SetupTest() {
	s.repo = &mocks.ITodoListRepository{}
	s.trxObj = &mocks.TrxObj{}
	s.usecase = todo_list_usecase.NewCrudTodoListUsecase(s.repo)
}

func TestCrudTodoListUsecase(t *testing.T) {
	suite.Run(t, new(CrudTodoListUsecaseTestSuite))
}

func (s *CrudTodoListUsecaseTestSuite) TestGetByUserID() {
	ctx := context.Background()
	userID := int64(1)

	testcases := []struct {
		name     string
		mockFunc func()
		wantLen  int
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				s.repo.On("GetByUserID", ctx, userID).Return([]*mentity.TodoList{
					{ID: 1, Title: "Test", UserID: userID},
				}, nil).Once()
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Error Repo",
			mockFunc: func() {
				s.repo.On("GetByUserID", ctx, userID).Return(nil, errors.New("db error")).Once()
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			res, err := s.usecase.GetByUserID(ctx, userID)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Len(res, tt.wantLen)
			}
		})
	}
}

func (s *CrudTodoListUsecaseTestSuite) TestGetByID() {
	ctx := context.Background()
	id := int64(1)

	testcases := []struct {
		name     string
		mockFunc func()
		wantNil  bool // want result to be nil (for not found case)
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				s.repo.On("GetByID", ctx, id).Return(&mentity.TodoList{ID: 1, Title: "Test"}, nil).Once()
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Not Found (Repo returns nil data, nil error)",
			// Assuming repository pattern might return nil, nil if using Take() with error handling
			// But wait, code says: if data == nil { return nil, nil }
			// However GORM usually returns ErrRecordNotFound.
			// usecase handles `if err != nil` then `if data == nil`.
			// If repo returns error, usecase returns error.
			// If repo returns nil data, nil error, usecase returns nil, nil.
			mockFunc: func() {
				s.repo.On("GetByID", ctx, id).Return(nil, nil).Once()
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Error Repo",
			mockFunc: func() {
				s.repo.On("GetByID", ctx, id).Return(nil, errors.New("db error")).Once()
			},
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			res, err := s.usecase.GetByID(ctx, id)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				if tt.wantNil {
					s.Nil(res)
				} else {
					s.NotNil(res)
				}
			}
		})
	}
}

func (s *CrudTodoListUsecaseTestSuite) TestCreate() {
	ctx := context.Background()
	req := entity.TodoListReq{
		UserID:      1,
		Title:       "Test",
		Description: "Desc",
		DoingAt:     "2023-10-27",
	}

	testcases := []struct {
		name     string
		req      entity.TodoListReq
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success",
			req:  req,
			mockFunc: func() {
				// usecase.ValidateStruct -> ok
				// ParseDate -> ok
				// repo.Create -> ok
				s.repo.On("Create", ctx, mock.Anything, mock.Anything, false).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "Validation Error",
			req:  entity.TodoListReq{}, // Empty title usually fails validation if tags exist
			mockFunc: func() {
				// ValidateStruct will fail. Repo not called.
			},
			wantErr: true,
		},
		{
			name: "Repo Error",
			req:  req,
			mockFunc: func() {
				s.repo.On("Create", ctx, mock.Anything, mock.Anything, false).Return(errors.New("db error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			res, err := s.usecase.Create(ctx, tt.req)
			if tt.wantErr {
				s.Error(err)
				s.Nil(res)
			} else {
				s.NoError(err)
				s.NotNil(res)
			}
		})
	}
}

func (s *CrudTodoListUsecaseTestSuite) TestUpdateByID() {
	ctx := context.Background()
	req := entity.TodoListReq{
		ID:          1,
		UserID:      1,
		Title:       "Updated",
		Description: "Updated Desc",
		DoingAt:     "2023-10-28",
	}

	testcases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				// DBTransaction sequence
				s.repo.On("Begin").Return(s.trxObj, nil).Once()

				// Inside callback:
				// LockByID
				s.repo.On("LockByID", ctx, s.trxObj, req.ID).
					Return(&mentity.TodoList{ID: 1, Title: "Old"}, nil).Once()

				// Update
				s.repo.On("Update", ctx, s.trxObj, mock.Anything, mock.Anything).
					Return(nil).Once()

				// Commit
				s.trxObj.On("Commit").Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "Begin Error",
			mockFunc: func() {
				s.repo.On("Begin").Return(nil, errors.New("begin error")).Once()
			},
			wantErr: true,
		},
		{
			name: "LockByID Error",
			mockFunc: func() {
				s.repo.On("Begin").Return(s.trxObj, nil).Once()
				s.repo.On("LockByID", ctx, s.trxObj, req.ID).
					Return(nil, errors.New("lock error")).Once()
				// Rollback triggered by defer
				s.trxObj.On("Rollback").Return(nil).Once()
			},
			wantErr: true,
		},
		{
			name: "Data Not Found (LockByID returns nil, nil)",
			mockFunc: func() {
				s.repo.On("Begin").Return(s.trxObj, nil).Once()
				s.repo.On("LockByID", ctx, s.trxObj, req.ID).
					Return(nil, nil).Once()
				s.trxObj.On("Rollback").Return(nil).Once()
			},
			wantErr: true,
		},
		{
			name: "Update Error",
			mockFunc: func() {
				s.repo.On("Begin").Return(s.trxObj, nil).Once()
				s.repo.On("LockByID", ctx, s.trxObj, req.ID).
					Return(&mentity.TodoList{ID: 1}, nil).Once()
				s.repo.On("Update", ctx, s.trxObj, mock.Anything, mock.Anything).
					Return(errors.New("update error")).Once()
				s.trxObj.On("Rollback").Return(nil).Once()
			},
			wantErr: true,
		},
		{
			name: "Commit Error",
			mockFunc: func() {
				s.repo.On("Begin").Return(s.trxObj, nil).Once()
				s.repo.On("LockByID", ctx, s.trxObj, req.ID).
					Return(&mentity.TodoList{ID: 1}, nil).Once()
				s.repo.On("Update", ctx, s.trxObj, mock.Anything, mock.Anything).
					Return(nil).Once()
				s.trxObj.On("Commit").Return(errors.New("commit error")).Once()
				// If commit fails, rollback might be called depending on implementation logic,
				// checking db transaction implementation:
				// if err = trx.Commit(); err != nil { return err } -> commit=true not reached.
				// defer checks !commit -> Rollback called.
				s.trxObj.On("Rollback").Return(nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			err := s.usecase.UpdateByID(ctx, req)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *CrudTodoListUsecaseTestSuite) TestDeleteByID() {
	ctx := context.Background()
	id := int64(1)

	testcases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				s.repo.On("DeleteByID", ctx, mock.Anything, id).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "Error",
			mockFunc: func() {
				s.repo.On("DeleteByID", ctx, mock.Anything, id).Return(errors.New("db error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range testcases {
		s.Run(tt.name, func() {
			tt.mockFunc()
			err := s.usecase.DeleteByID(ctx, id)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
