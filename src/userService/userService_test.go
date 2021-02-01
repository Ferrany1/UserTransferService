package userService

import (
	"UserTransferService/src/db_service"
	"UserTransferService/src/system/config"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var (
	testUser1 = UserInfo{UserData: db_service.User{Username: "test1", Email: "test1@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test1", Balance: 10}}
	updUser   = UserInfo{UserData: db_service.User{Username: "test1", Email: "test2@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test1", Balance: 10}}
	receiver  = UserInfo{UserData: db_service.User{Username: "test2", Email: "test2@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test2", Balance: 5}}
	tr        = 5
)

type HandlerTest struct {
	suite.Suite
	mdb  *gorm.DB
	mock sqlmock.Sqlmock
}

func (m *HandlerTest) AfterTest() {
	require.NoError(m.T(), m.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(HandlerTest))
}

func (m *HandlerTest) SetupSuite() { // or *gorm.DB
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

	db_service.DB.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(m.T(), err)
}

func (m *HandlerTest) Test_verifyUser() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	_, err := testUser1.verifyUser()
	require.NoError(m.T(), err)
}

func (m *HandlerTest) TestUserInfo_CreateNewUser() {
	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`INSERT INTO "users" ("username","email","password") VALUES ($1,$2,$3)`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`INSERT INTO "balances" ("username","balance") VALUES ($1,$2)`).
		WithArgs(testUser1.UserBalance.Username, 0).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := testUser1.CreateNewUser()
	require.NoError(m.T(), err)
}

func (m *HandlerTest) TestUserInfo_GetUserBalance() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectQuery(
		`SELECT * FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(testUser1.UserBalance.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "balance"}).
			AddRow(testUser1.UserBalance.Username, testUser1.UserBalance.Balance))

	b, err := testUser1.GetUserBalance()
	require.Equal(m.T(), b, db_service.Balance{Username: testUser1.UserBalance.Username, Balance: testUser1.UserBalance.Balance})
	require.NoError(m.T(), err)
}

func (m *HandlerTest) TestUserInfo_UpdateUser() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`UPDATE "users" SET "username"=$1,"email"=$2,"password"=$3 WHERE "users"."username" = $4`).
		WithArgs(testUser1.UserData.Username, updUser.UserData.Email, testUser1.UserData.Password, testUser1.UserData.Username).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	err := testUser1.UpdateUser(updUser)
	require.NoError(m.T(), err)
}

func (m *HandlerTest) TestUserInfo_TransferUsersBalance() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(receiver.UserData.Username, receiver.UserData.Email, receiver.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(receiver.UserData.Username, receiver.UserData.Email, receiver.UserData.Password))

	m.mock.ExpectQuery(
		`SELECT * FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(testUser1.UserBalance.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "balance"}).
			AddRow(testUser1.UserBalance.Username, testUser1.UserBalance.Balance))

	m.mock.ExpectQuery(
		`SELECT * FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(receiver.UserBalance.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "balance"}).
			AddRow(receiver.UserBalance.Username, receiver.UserBalance.Balance))

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`UPDATE "balances" SET "username"=$1,"balance"=$2 WHERE "balances"."username" = $3`).
		WithArgs(testUser1.UserBalance.Username, testUser1.UserBalance.Balance - tr, testUser1.UserBalance.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))

	m.mock.ExpectExec("SAVEPOINT sp0x16c7640").WillReturnResult(sqlmock.NewResult(1, 1))

	m.mock.ExpectExec(
		`UPDATE "balances" SET "username"=$1,"balance"=$2 WHERE "balances"."username" = $3`).
		WithArgs(receiver.UserBalance.Username, receiver.UserBalance.Balance + tr, receiver.UserBalance.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := testUser1.TransferUsersBalance(receiver, tr)
	require.NoError(m.T(), err)
}

func (m *HandlerTest) TestUserInfo_DeleteUser() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`DELETE FROM "users" WHERE "users"."username" = $1`).
		WithArgs(testUser1.UserData.Username).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`DELETE FROM "balances" WHERE "balances"."username" = $1`).
		WithArgs(testUser1.UserData.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()

	err := testUser1.DeleteUser()
	require.NoError(m.T(), err)
}
