package db_service

import (
	"UserTransferService/src/system/config"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

var (
	testUserData1    = User{Username: "test1", Email: "test1@test.com", Password: "safePassword"}
	testUserBalance1 = Balance{Username: "test1", Balance: 10}
	testUserBalance2 = Balance{Username: "test2", Balance: 20}
)

type DBTest struct {
	suite.Suite
	mdb  Database
	mock sqlmock.Sqlmock
}

func (m *DBTest) AfterTest() {
	require.NoError(m.T(), m.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(DBTest))
}

func (m *DBTest) SetupSuite() { // or *gorm.DB
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(m.T(), err)

	config.CF = new(config.Config)
	config.CF.DB.Tables.Users = "users"
	config.CF.DB.Tables.Balances = "balances"

	m.mock = mock
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "db_service",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	m.mdb.DB, err = gorm.Open(dialector, &gorm.Config{})
	m.mdb.DB.Logger.LogMode(logger.Silent)
	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_GetUser() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUserData1.Username, testUserData1.Email, testUserData1.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUserData1.Username, testUserData1.Email, testUserData1.Password))

	u, err := m.mdb.GetUser(testUserData1)

	require.Equal(m.T(), testUserData1, u)
	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_GetBalance() {
	m.mock.ExpectQuery(
		`SELECT * FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(testUserBalance1.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "balance"}).
			AddRow(testUserBalance1.Username, testUserBalance1.Balance))

	b, err := m.mdb.GetBalance(testUserData1)

	require.Equal(m.T(), testUserBalance1, b)
	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_CreateUser() {
	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`INSERT INTO "users" ("username","email","password") VALUES ($1,$2,$3)`).
		WithArgs(testUserData1.Username, testUserData1.Email, testUserData1.Password).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`INSERT INTO "balances" ("username","balance") VALUES ($1,$2)`).
		WithArgs(testUserBalance1.Username, 0).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := m.mdb.CreateUser(testUserData1)
	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_UpdateUser() {
	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`UPDATE "users" SET "username"=$1,"email"=$2,"password"=$3 WHERE "users"."username" = $4`).
		WithArgs(testUserData1.Username, testUserData1.Email, testUserData1.Password, testUserData1.Username).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	err := m.mdb.UpdateUser(testUserData1)

	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_UpdateBalance() {
	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`UPDATE "balances" SET "username"=$1,"balance"=$2 WHERE "balances"."username" = $3`).
		WithArgs(testUserBalance1.Username, testUserBalance1.Balance, testUserBalance1.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))

	m.mock.ExpectExec("SAVEPOINT sp0x174a5c0").WillReturnResult(sqlmock.NewResult(1, 1))

	m.mock.ExpectExec(
		`UPDATE "balances" SET "username"=$1,"balance"=$2 WHERE "balances"."username" = $3`).
		WithArgs(testUserBalance2.Username, testUserBalance2.Balance, testUserBalance2.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := m.mdb.UpdateBalance(testUserBalance1, testUserBalance2)

	require.NoError(m.T(), err)
}

func (m *DBTest) TestDatabase_DeleteUser() {
	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`DELETE FROM "users" WHERE "users"."username" = $1`).
		WithArgs(testUserData1.Username).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`DELETE FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(testUserBalance1.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := m.mdb.DeleteUser(testUserData1)

	require.NoError(m.T(), err)
}
