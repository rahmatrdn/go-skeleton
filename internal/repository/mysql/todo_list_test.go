package mysql_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	errwrap "github.com/pkg/errors"
	"github.com/rahmatrdn/go-skeleton/config"
	apperr "github.com/rahmatrdn/go-skeleton/error"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"
	"github.com/stretchr/testify/suite"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoListRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
	repo *mysql.TodoListRepository
}

func TestTodoListRepository(t *testing.T) {
	suite.Run(t, new(TodoListRepositoryTestSuite))
}

func (s *TodoListRepositoryTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	dialector := gmysql.New(gmysql.Config{
		Conn:                      s.db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	s.Require().NoError(err)

	s.repo = mysql.NewTodoListRepository(&config.Mysql{DB: gormDB})
}

func (s *TodoListRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *TodoListRepositoryTestSuite) TestGetByUserID() {
	type args struct {
		ctx    context.Context
		userID int64
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately for context cancelled case

	tests := []struct {
		name         string
		args         args
		mockSetup    func()
		wantLen      int
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mockSetup: func() {
				expectedRows := sqlmock.NewRows([]string{"id", "user_id", "title"}).
					AddRow(1, 1, "Test Todo")
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE user_id = ?")).
					WithArgs(1).
					WillReturnRows(expectedRows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Error Query",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE user_id = ?")).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantLen: 0,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(sql.ErrConnDone, err)
			},
		},
		{
			name: "Not Found",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE user_id = ?")).
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantLen: 0,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrRecordNotFound(), err)
			},
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:    cancelledCtx,
				userID: 1,
			},
			mockSetup: func() {
				// No expectations because CheckDeadline returns error first
			},
			wantLen: 0,
			wantErr: true,
			inspectError: func(err error) {
				s.True(errwrap.Is(err, context.Canceled))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			result, err := s.repo.GetByUserID(tt.args.ctx, tt.args.userID)
			if tt.wantErr {
				s.Error(err)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				s.Len(result, tt.wantLen)
			}
		})
	}
}

func (s *TodoListRepositoryTestSuite) TestGetByID() {
	type args struct {
		ctx context.Context
		id  int64
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name         string
		args         args
		mockSetup    func()
		want         *entity.TodoList
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockSetup: func() {
				expectedRows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "Test Todo")
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE id = ? LIMIT 1")).
					WithArgs(1).
					WillReturnRows(expectedRows)
			},
			want: &entity.TodoList{
				ID:    1,
				Title: "Test Todo",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE id = ? LIMIT 1")).
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrRecordNotFound(), err)
			},
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx: cancelledCtx,
				id:  1,
			},
			mockSetup: func() {},
			want:      nil,
			wantErr:   true,
			inspectError: func(err error) {
				s.True(errwrap.Is(err, context.Canceled))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			result, err := s.repo.GetByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				s.Error(err)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				s.Equal(tt.want.Title, result.Title)
			}
		})
	}
}

func (s *TodoListRepositoryTestSuite) TestCreate() {
	type args struct {
		ctx        context.Context
		dbTrx      mysql.TrxObj
		params     *entity.TodoList
		nonZeroVal bool
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				ctx:        context.Background(),
				dbTrx:      new(mocks.TrxObj),
				params:     &entity.TodoList{Title: "New Todo", UserID: 1},
				nonZeroVal: true,
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `todo_lists`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error DB",
			args: args{
				ctx:        context.Background(),
				dbTrx:      new(mocks.TrxObj),
				params:     &entity.TodoList{Title: "New Todo", UserID: 1},
				nonZeroVal: true,
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `todo_lists`")).
					WillReturnError(sql.ErrConnDone)
				s.mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:        cancelledCtx,
				dbTrx:      new(mocks.TrxObj),
				params:     &entity.TodoList{Title: "New Todo", UserID: 1},
				nonZeroVal: true,
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			err := s.repo.Create(tt.args.ctx, tt.args.dbTrx, tt.args.params, tt.args.nonZeroVal)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *TodoListRepositoryTestSuite) TestLockByID() {
	type args struct {
		ctx   context.Context
		dbTrx mysql.TrxObj
		ID    int64
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name         string
		args         args
		mockSetup    func()
		want         *entity.TodoList
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				ID:    1,
			},
			mockSetup: func() {
				expectedRows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow(1, "Locked Todo")
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE id = ? FOR UPDATE")).
					WithArgs(1).
					WillReturnRows(expectedRows)
			},
			want: &entity.TodoList{
				ID:    1,
				Title: "Locked Todo",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				ID:    1,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todo_lists WHERE id = ? FOR UPDATE")).
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrRecordNotFound(), err)
			},
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:   cancelledCtx,
				dbTrx: new(mocks.TrxObj),
				ID:    1,
			},
			mockSetup: func() {},
			want:      nil,
			wantErr:   true,
			inspectError: func(err error) {
				s.True(errwrap.Is(err, context.Canceled))
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			result, err := s.repo.LockByID(tt.args.ctx, tt.args.dbTrx, tt.args.ID)
			if tt.wantErr {
				s.Error(err)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.want.Title, result.Title)
			}
		})
	}
}

func (s *TodoListRepositoryTestSuite) TestUpdate() {
	type args struct {
		ctx     context.Context
		dbTrx   mysql.TrxObj
		params  *entity.TodoList
		changes *entity.TodoList
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				params: &entity.TodoList{
					ID:    1,
					Title: "Old Title",
				},
				changes: &entity.TodoList{
					Title: "New Title",
				},
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `todo_lists`")).
					WillReturnResult(sqlmock.NewResult(0, 1))
				s.mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Success No Changes",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				params: &entity.TodoList{
					ID:    1,
					Title: "Old Title",
				},
				changes: nil,
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `todo_lists`")).
					WillReturnResult(sqlmock.NewResult(0, 1))
				s.mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error DB",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				params: &entity.TodoList{
					ID:    1,
					Title: "Old Title",
				},
				changes: &entity.TodoList{
					Title: "New Title",
				},
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `todo_lists`")).
					WillReturnError(sql.ErrConnDone)
				s.mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:   cancelledCtx,
				dbTrx: new(mocks.TrxObj),
				params: &entity.TodoList{
					ID:    1,
					Title: "Old Title",
				},
				changes: &entity.TodoList{
					Title: "New Title",
				},
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			err := s.repo.Update(tt.args.ctx, tt.args.dbTrx, tt.args.params, tt.args.changes)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *TodoListRepositoryTestSuite) TestDeleteByID() {
	type args struct {
		ctx   context.Context
		dbTrx mysql.TrxObj
		id    int64
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				id:    123,
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `todo_lists` WHERE id = ?")).
					WithArgs(123).
					WillReturnResult(sqlmock.NewResult(0, 1))
				s.mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error DB",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				id:    123,
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `todo_lists` WHERE id = ?")).
					WithArgs(123).
					WillReturnError(sql.ErrConnDone)
				s.mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:   cancelledCtx,
				dbTrx: new(mocks.TrxObj),
				id:    123,
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			err := s.repo.DeleteByID(tt.args.ctx, tt.args.dbTrx, tt.args.id)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
