package mysql_test

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"
	"github.com/rahmatrdn/go-skeleton/tests/fixture"
	"github.com/rahmatrdn/go-skeleton/tests/fixture/factory"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	walletColumns = factory.GetCols(&entity.Wallet{})
)

type WalletTestSuite struct {
	suite.Suite
}

func (s *WalletTestSuite) SetupTest() {
}

func TestWalletTest(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) TestGetByID() {
	wallet := factory.StubbedWallet()

	testcases := map[string]struct {
		ctx      context.Context
		result   *entity.Wallet
		getError error
		wantErr  bool
	}{
		"fail ctx ended": {
			ctx:      fixture.CtxEnded(),
			result:   wallet,
			getError: nil,
			wantErr:  true,
		},
		"fail get not found": {
			ctx:      context.Background(),
			getError: gorm.ErrRecordNotFound,
			wantErr:  true,
		},
		"fail get from db": {
			ctx:      context.Background(),
			getError: errors.New("fail get"),
			wantErr:  true,
		},
		"success": {
			ctx:      context.Background(),
			result:   wallet,
			getError: nil,
			wantErr:  false,
		},
	}

	for name, tc := range testcases {
		db, mock, err := sqlmock.New()
		if err != nil {
			s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		if tc.getError != nil {
			mock.ExpectQuery("^SELECT(.+)").WillReturnError(tc.getError)
		} else {
			rows := sqlmock.NewRows(walletColumns)
			if tc.result != nil {
				rows.AddRow(factory.GetRows(wallet, false)...)
			}
			mock.ExpectQuery("^SELECT(.+)").WillReturnRows(rows)
		}

		gormDB, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db, SkipInitializeWithVersion: true}))
		if err != nil {
			panic(err)
		}

		repo := mysql.NewWalletRepository(&config.Mysql{DB: gormDB})
		res, err := repo.GetByID(tc.ctx, wallet.ID)

		if tc.wantErr {
			s.Error(err, name)
		} else {
			s.NoError(err, name)
			s.Equal(tc.result, res)
		}
	}
}

func (s *WalletTestSuite) TestCreate() {
	wallet := factory.StubbedWallet()

	testcases := map[string]struct {
		ctx        context.Context
		model      *entity.Wallet
		nonZeroVal bool
		expErr     error
	}{
		"success": {
			ctx:        context.Background(),
			model:      wallet,
			nonZeroVal: true,
		},
		"error": {
			ctx:        context.Background(),
			model:      wallet,
			expErr:     errors.New(uuid.NewString()),
			nonZeroVal: true,
		},
		"fail ctx ended": {
			ctx:        fixture.CtxEnded(),
			model:      wallet,
			expErr:     errors.Wrap(context.Canceled, "WalletRepository.Create"),
			nonZeroVal: true,
		},
	}

	for name, tc := range testcases {
		tc := tc
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()

			if tc.expErr != nil {
				mock.ExpectExec("^INSERT INTO(.+)").WillReturnError(tc.expErr)
				mock.ExpectRollback()
			} else {
				rows := factory.GetInsertRows(tc.model, tc.nonZeroVal, map[string]any{"ID": tc.model.ID})
				mock.ExpectExec("^INSERT INTO(.+)").WithArgs(rows...).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			gormDB, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db, SkipInitializeWithVersion: true}))
			if err != nil {
				panic(err)
			}

			repo := mysql.NewWalletRepository(&config.Mysql{DB: gormDB})
			err = repo.Create(tc.ctx, nil, tc.model, tc.nonZeroVal)
			if tc.expErr != nil {
				if s.Error(err, name) {
					s.Equal(errors.Cause(tc.expErr), errors.Cause(err))
				}
			} else {
				s.NoError(err, name)
			}
		})
	}
}

func (s *WalletTestSuite) TestLockByID() {
	wallet := factory.StubbedWallet()

	testcases := map[string]struct {
		ctx      context.Context
		result   *entity.Wallet
		getError error
		wantErr  bool
	}{
		"fail ctx ended": {
			ctx:      fixture.CtxEnded(),
			result:   nil,
			getError: nil,
			wantErr:  true,
		},
		"fail get not found": {
			ctx:      context.Background(),
			getError: gorm.ErrRecordNotFound,
			wantErr:  true,
		},
		"fail get error from db": {
			ctx:      context.Background(),
			getError: errors.New("fail get"),
			wantErr:  true,
		},
		"success": {
			ctx:      context.Background(),
			result:   wallet,
			getError: nil,
			wantErr:  false,
		},
	}

	for name, tc := range testcases {
		tc := tc
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tc.getError != nil {
				mock.ExpectQuery("^SELECT(.+)WHERE id =(.+)FOR UPDATE").WillReturnError(tc.getError)
			} else {
				rows := sqlmock.NewRows(walletColumns)
				if tc.result != nil {
					rows.AddRow(factory.GetRows(wallet, false)...)
				}
				mock.ExpectQuery("^SELECT(.+)WHERE id =(.+)FOR UPDATE").WillReturnRows(rows)
			}

			gormDB, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db, SkipInitializeWithVersion: true}))
			if err != nil {
				panic(err)
			}

			repo := mysql.NewWalletRepository(&config.Mysql{DB: gormDB})
			res, err := repo.LockByID(tc.ctx, nil, wallet.ID)

			if tc.wantErr {
				s.Error(err, name)
			} else {
				if s.NoError(err, name) {
					s.Equal(tc.result, res)
				}
			}
		})
	}
}

func (s *WalletTestSuite) TestUpdate() {
	testcases := map[string]struct {
		ctx     context.Context
		model   *entity.Wallet
		changes *entity.Wallet
		expErr  error
	}{
		"success": {
			ctx:   context.Background(),
			model: factory.StubbedWallet(),
		},
		// "with_changes": {
		// 	ctx:   context.Background(),
		// 	model: factory.StubbedWallet(),
		// 	changes: &entity.Wallet{
		// 		UserName: "Update Name",
		// 	},
		// },
		"error": {
			ctx:    context.Background(),
			model:  factory.StubbedWallet(),
			expErr: errors.New(uuid.NewString()),
		},
		"fail ctx ended": {
			ctx:    fixture.CtxEnded(),
			model:  factory.StubbedWallet(),
			expErr: errors.Wrap(context.Canceled, "WalletRepository.Update"),
		},
	}

	for name, tc := range testcases {
		name, tc := name, tc
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()

			if tc.expErr != nil {
				mock.ExpectExec("^UPDATE `wallets`(.+)").WillReturnError(tc.expErr)
				mock.ExpectRollback()
			} else {
				var rows []driver.Value
				if tc.changes != nil {
					rows = factory.GetUpdateRows(tc.changes, true, "ID", tc.model.ID, true)
				} else {
					rows = factory.GetUpdateRows(tc.model, false, "ID", tc.model.ID, true)
				}
				mock.ExpectExec("^UPDATE `wallets`(.+)").WithArgs(rows...).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			gormDB, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db, SkipInitializeWithVersion: true}))
			if err != nil {
				panic(err)
			}

			repo := mysql.NewWalletRepository(&config.Mysql{DB: gormDB})
			err = repo.Update(tc.ctx, nil, tc.model, tc.changes)

			if tc.expErr != nil {
				if s.Error(err, name) {
					s.Equal(errors.Cause(tc.expErr), errors.Cause(err))
				}
			} else {
				s.NoError(err, name)
			}
		})
	}
}

func (s *WalletTestSuite) TestDeleteByID() {
	testcases := map[string]struct {
		ctx    context.Context
		model  *entity.Wallet
		expErr error
	}{
		"success": {
			ctx:   context.Background(),
			model: &entity.Wallet{ID: 1},
		},
		"error": {
			ctx:    context.Background(),
			model:  &entity.Wallet{ID: 2},
			expErr: errors.New(uuid.NewString()),
		},
		"fail ctx ended": {
			ctx:    fixture.CtxEnded(),
			model:  &entity.Wallet{ID: 3},
			expErr: errors.Wrap(context.Canceled, "WalletRepository.DeleteByID"),
		},
	}

	for name, tc := range testcases {
		name, tc := name, tc
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()

			if tc.expErr != nil {
				mock.ExpectExec("^DELETE FROM (.+)").WillReturnError(tc.expErr)
				mock.ExpectRollback()
			} else {
				mock.ExpectExec("^DELETE FROM (.+)").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			gormDB, err := gorm.Open(gmysql.New(gmysql.Config{Conn: db, SkipInitializeWithVersion: true}))
			if err != nil {
				panic(err)
			}

			repo := mysql.NewWalletRepository(&config.Mysql{DB: gormDB})
			err = repo.DeleteByID(tc.ctx, nil, tc.model.ID)

			if tc.expErr != nil {
				if s.Error(err, name) {
					s.Equal(errors.Cause(tc.expErr), errors.Cause(err))
				}
			} else {
				s.NoError(err, name)
			}
		})
	}
}
