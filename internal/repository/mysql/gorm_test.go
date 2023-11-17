package mysql_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rahmatrdn/go-skeleton/config"
	"github.com/rahmatrdn/go-skeleton/tests/mocks"

	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql"
	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormTrxSupportTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
	repo *mysql.User
	d    *GormTrxSupportTestData
}

type GormTrxSupportTestData struct {
	ctx    context.Context
	cancel context.CancelFunc
	rows   *sqlmock.Rows
}

func TestGormTrxSupport(t *testing.T) {
	suite.Run(t, new(GormTrxSupportTestSuite))
}

func (s *GormTrxSupportTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *GormTrxSupportTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	dialector := gmysql.New(gmysql.Config{Conn: s.db, SkipInitializeWithVersion: true})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	s.repo = mysql.NewUserRepository(&config.Mysql{DB: gormDB})
	s.d = &GormTrxSupportTestData{}
}

func (s *GormTrxSupportTestSuite) setupValidData() {
	s.d.ctx, s.d.cancel = context.WithDeadline(context.Background(), time.Now().Add(time.Hour))
	fakeData := &entity.User{}
	if err := faker.FakeData(fakeData); err != nil {
		s.FailNowf("an error '%s' was no expected when trying to create fake data", err.Error())
	}
	rows := sqlmock.NewRows([]string{"count"})
	s.d.rows = rows.AddRow(1)
}

func (s *GormTrxSupportTestSuite) TestCommitSuccess() {
	s.setupValidData()

	s.mock.ExpectBegin()
	s.mock.ExpectCommit()

	trx, err := s.repo.Begin()
	s.NoError(err)
	s.NotNil(trx)

	err = trx.Commit()
	s.NoError(err)
}

func (s *GormTrxSupportTestSuite) TestRollbackSuccess() {
	s.setupValidData()

	s.mock.ExpectBegin()
	s.mock.ExpectRollback()

	trx, err := s.repo.Begin()
	s.NoError(err)
	s.NotNil(trx)

	err = trx.Rollback()
	s.NoError(err)
}

func (s *GormTrxSupportTestSuite) TestTrx() {
	trx, _ := s.repo.Begin()

	db := s.repo.Trx(trx)
	s.NotNil(db)

	db = s.repo.Trx(nil)
	s.NotNil(db)
}

// DB Transaction Test
type DBTransactionTestSuite struct {
	suite.Suite
	repo *mocks.UserRepository
	trx  *mocks.TrxObj
}

func TestDBTransaction(t *testing.T) {
	suite.Run(t, new(DBTransactionTestSuite))
}

func (s *DBTransactionTestSuite) SetupTest() {
	s.repo = &mocks.UserRepository{}
	s.trx = &mocks.TrxObj{}
}

func (s *DBTransactionTestSuite) setupValidMock() {
	s.repo.On("Begin").Return(s.trx, nil).Once()
	s.trx.On("Rollback").Return(nil).Once()
	s.trx.On("Commit").Return(nil).Once()
}

func (s *DBTransactionTestSuite) TestSuccess() {
	s.setupValidMock()
	err := mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { return nil })

	s.NoError(err)

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertCalled(s.T(), "Commit")
	s.trx.AssertNotCalled(s.T(), "Rollback")
}

func (s *DBTransactionTestSuite) TestErrorOnBegin() {
	randErr := errors.New(uuid.New().String())
	s.repo.On("Begin").Return(nil, randErr).Once()
	s.setupValidMock()
	err := mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { return nil })

	s.Error(err)
	s.Equal(randErr, errors.Cause(err))

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertNotCalled(s.T(), "Commit")
	s.trx.AssertNotCalled(s.T(), "Rollback")
}

func (s *DBTransactionTestSuite) TestErrorCallback() {
	randErr := errors.New(uuid.New().String())
	s.setupValidMock()
	err := mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { return randErr })

	s.Error(err)
	s.Equal(randErr, errors.Cause(err))

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertNotCalled(s.T(), "Commit")
	s.trx.AssertCalled(s.T(), "Rollback")
}

func (s *DBTransactionTestSuite) TestErrorOnCommit() {
	randErr := errors.New(uuid.New().String())
	s.trx.On("Commit").Return(randErr).Once()
	s.setupValidMock()
	err := mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { return nil })

	s.Error(err)
	s.Equal(randErr, errors.Cause(err))

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertCalled(s.T(), "Commit")
	s.trx.AssertCalled(s.T(), "Rollback")
}

func (s *DBTransactionTestSuite) TestErrorCallbackAndRollback() {
	cErr := errors.New(uuid.New().String())
	rErr := errors.New(uuid.New().String())
	s.trx.On("Rollback").Return(rErr).Once()
	s.setupValidMock()
	err := mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { return cErr })

	s.Error(err)
	s.Equal(rErr, errors.Cause(err))
	s.Contains(err.Error(), rErr.Error())
	s.Contains(err.Error(), cErr.Error())

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertNotCalled(s.T(), "Commit")
	s.trx.AssertCalled(s.T(), "Rollback")
}

func (s *DBTransactionTestSuite) TestPanicAndRollback() {
	cErr := errors.New(uuid.New().String())
	rErr := errors.New(uuid.New().String())
	s.trx.On("Rollback").Return(rErr).Once()
	s.setupValidMock()

	s.Panics(func() {
		mysql.DBTransaction(s.repo, func(trx mysql.TrxObj) error { panic(cErr) })
	})

	s.repo.AssertCalled(s.T(), "Begin")
	s.trx.AssertNotCalled(s.T(), "Commit")
	s.trx.AssertCalled(s.T(), "Rollback")
}
