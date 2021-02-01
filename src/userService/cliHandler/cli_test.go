package cliHandler

import (
	"UserTransferService/src/db_service"
	"UserTransferService/src/system/config"
	"UserTransferService/src/userService"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

var (
	testUser1 = userService.UserInfo{UserData: db_service.User{Username: "test1", Email: "test1@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test1", Balance: 10}}
	updUser   = userService.UserInfo{UserData: db_service.User{Username: "test1", Email: "test2@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test1", Balance: 10}}
	receiver  = userService.UserInfo{UserData: db_service.User{Username: "test2", Email: "test2@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test2", Balance: 5}}
	tr        = 5
)

type CliTest struct {
	suite.Suite
	mdb  	*gorm.DB
	mock 	sqlmock.Sqlmock
	ca 		consoleArgs
}

func (m *CliTest) AfterTest() {
	require.NoError(m.T(), m.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(CliTest))
}

func (m *CliTest) SetupSuite() { // or *gorm.D

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

func (m *CliTest) TestCreateUser() {
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

	testText := fmt.Sprintf("create username:%s email:%s password:%s", testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommandCheck(testText)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	require.Equal(m.T(), "> Success: created new user\n", string(out))
}

func (m *CliTest) TestGetUserBalance() {
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

	testText := fmt.Sprintf("balance username:%s email:%s password:%s", testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommandCheck(testText)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	require.Equal(m.T(), "> Balance: 10\n", string(out))
}

func (m *CliTest) TestTransferFromBalance() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1`).
		WithArgs(receiver.UserData.Username).
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


	testText := fmt.Sprintf("transfer username:%s email:%s password:%s receiver:%s amount:%s", testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password, receiver.UserBalance.Username, strconv.Itoa(tr))

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommandCheck(testText)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	require.Equal(m.T(), "> Success: made a transaction\n", string(out))
}

func (m *CliTest) TestUpdateUser() {
	m.mock.ExpectQuery(
		`SELECT * FROM "users" WHERE "users"."username" = $1 AND "users"."email" = $2 AND "users"."password" = $3`).
		WithArgs(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password).
		WillReturnRows(sqlmock.NewRows([]string{"username", "email", "password"}).
			AddRow(testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password))

	m.mock.ExpectBegin()
	m.mock.ExpectExec(
		`UPDATE "users" SET "username"=$1,"email"=$2,"password"=$3 WHERE "users"."username" = $4`).
		WithArgs(testUser1.UserData.Username, "newEmail", "newPassword", testUser1.UserData.Username).
		WillReturnResult(sqlmock.NewResult(1, 3))
	m.mock.ExpectCommit()

	testText := fmt.Sprintf("update username:%s email:%s password:%s new_email:%s new_password:%s", testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password, "newEmail", "newPassword")

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommandCheck(testText)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	require.Equal(m.T(), "> Success: updated user info\n", string(out))
}

func (m *CliTest) TestDeleteUser() {
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

	testText := fmt.Sprintf("delete username:%s email:%s password:%s", testUser1.UserData.Username, testUser1.UserData.Email, testUser1.UserData.Password)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	CommandCheck(testText)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	require.Equal(m.T(), "> Success: deleted user\n", string(out))
}