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

type UserRepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
	repo *mysql.User
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	dialector := gmysql.New(gmysql.Config{
		Conn:                      s.db,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	s.Require().NoError(err)

	s.repo = mysql.NewUserRepository(&config.Mysql{DB: gormDB})
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *UserRepositoryTestSuite) TestCreate() {
	type args struct {
		ctx   context.Context
		dbTrx mysql.TrxObj
		user  *entity.User
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
				user: &entity.User{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error DB",
			args: args{
				ctx:   context.Background(),
				dbTrx: new(mocks.TrxObj),
				user: &entity.User{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockSetup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
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
				user: &entity.User{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockSetup: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockSetup()
			err := s.repo.Create(tt.args.ctx, tt.args.dbTrx, tt.args.user)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestLockByID() {
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
		want         *entity.User
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
				expectedRows := sqlmock.NewRows([]string{"id", "email", "name"}).
					AddRow(1, "test@example.com", "Test User")
				// GORM v2 parameterizes LIMIT: LIMIT ?
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? LIMIT ? FOR UPDATE")).
					WithArgs(1, 1).
					WillReturnRows(expectedRows)
			},
			want: &entity.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
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
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? LIMIT ? FOR UPDATE")).
					WithArgs(1, 1). // ID, Limit
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrUserNotFound(), err)
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
				s.Equal(tt.want.Email, result.Email)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetByEmail() {
	type args struct {
		ctx   context.Context
		email string
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name         string
		args         args
		mockSetup    func()
		want         *entity.User
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			mockSetup: func() {
				expectedRows := sqlmock.NewRows([]string{"id", "email", "name"}).
					AddRow(1, "test@example.com", "Test User")
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? LIMIT ?")).
					WithArgs("test@example.com", 1).
					WillReturnRows(expectedRows)
			},
			want: &entity.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{
				ctx:   context.Background(),
				email: "notfound@example.com",
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? LIMIT ?")).
					WithArgs("notfound@example.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrUserNotFound(), err)
			},
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:   cancelledCtx,
				email: "test@example.com",
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
			result, err := s.repo.GetByEmail(tt.args.ctx, tt.args.email)
			if tt.wantErr {
				s.Error(err)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.want.Email, result.Email)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetByEmailAndRole() {
	type args struct {
		ctx   context.Context
		email string
		role  int8
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name         string
		args         args
		mockSetup    func()
		want         *entity.User
		wantErr      bool
		inspectError func(err error)
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
				role:  1,
			},
			mockSetup: func() {
				expectedRows := sqlmock.NewRows([]string{"id", "email", "name", "role"}).
					AddRow(1, "test@example.com", "Test User", 1)
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND role = ? LIMIT ?")).
					WithArgs("test@example.com", int8(1), 1).
					WillReturnRows(expectedRows)
			},
			want: &entity.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
				Role:  1,
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{
				ctx:   context.Background(),
				email: "notfound@example.com",
				role:  1,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND role = ? LIMIT ?")).
					WithArgs("notfound@example.com", int8(1), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
			inspectError: func(err error) {
				s.Equal(apperr.ErrUserNotFound(), err)
			},
		},
		{
			name: "Context Cancelled",
			args: args{
				ctx:   cancelledCtx,
				email: "test@example.com",
				role:  1,
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
			// Cast role to entity.RoleType if needed, but entity definition uses int8 in struct
			// though I saw entity.RoleType in repository interface method signature.
			// Let's verify entity.RoleType from previous view_file or list?
			// Wait, I saw "role entity.RoleType" in user.go line 71.
			// However in entity/user.go line 9 it is "Role int8".
			// I should probably check entity definition for RoleType, assuming it is int8alias or similar.
			// Or just pass int8(tt.args.role) if it complains, but let's assume it matches for now.

			// Actually, let me check entity.Role type if possible or just use int8 in tests struct and cast.
			// The user.go file says: `role entity.RoleType`
			// I must cast it.

			result, err := s.repo.GetByEmailAndRole(tt.args.ctx, tt.args.email, entity.RoleType(tt.args.role))
			if tt.wantErr {
				s.Error(err)
				if tt.inspectError != nil {
					tt.inspectError(err)
				}
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.want.Email, result.Email)
			}
		})
	}
}
