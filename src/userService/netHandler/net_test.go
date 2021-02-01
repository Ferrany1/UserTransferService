package netHandler

import (
	"UserTransferService/src/db_service"
	"UserTransferService/src/system/config"
	"UserTransferService/src/userService"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testUser1 = userService.UserInfo{UserData: db_service.User{Username: "test1", Email: "test1@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test1", Balance: 10}}
	receiver  = userService.UserInfo{UserData: db_service.User{Username: "test2", Email: "test2@test.com", Password: "safePassword"}, UserBalance: db_service.Balance{Username: "test2", Balance: 5}}
	tr        = 5
)

type NetTest struct {
	suite.Suite
	mdb  	*gorm.DB
	mock 	sqlmock.Sqlmock
	router 	*gin.Engine
}

func (m *NetTest) AfterTest() {
	require.NoError(m.T(), m.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(NetTest))
}

func (m *NetTest) SetupSuite() { // or *gorm.D
	m.createRouter()

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

func (m *NetTest) TestCreateUser() {
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

	b, err := json.Marshal(testUser1)
	require.NoError(m.T(), err)

	w := m.performRequest(http.MethodPost, "/api/v1/user/create", nil, bytes.NewReader(b))

	b, err = ioutil.ReadAll(w.Body)
	require.NoError(m.T(), err)

	require.Equal(m.T(), http.StatusOK, w.Code)
}

func (m *NetTest) TestGetUserBalance() {
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

	b, err := json.Marshal(testUser1)
	require.NoError(m.T(), err)

	w := m.performRequest(http.MethodGet, "/api/v1/user/balance", nil, bytes.NewReader(b))

	b, err = ioutil.ReadAll(w.Body)
	require.NoError(m.T(), err)

	require.Equal(m.T(), http.StatusOK, w.Code)
}

func (m *NetTest) TestTransferFromBalance() {
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

	m.mock.ExpectExec("SAVEPOINT sp0x1786060").WillReturnResult(sqlmock.NewResult(1, 1))

	m.mock.ExpectExec(
		`UPDATE "balances" SET "username"=$1,"balance"=$2 WHERE "balances"."username" = $3`).
		WithArgs(receiver.UserBalance.Username, receiver.UserBalance.Balance + tr, receiver.UserBalance.Username).
		WillReturnResult(sqlmock.NewResult(1, 2))
	m.mock.ExpectCommit()


	b, err := json.Marshal(testUser1)
	require.NoError(m.T(), err)

	w := m.performRequest(http.MethodPost, "/api/v1/user/transfer/test2/5", nil, bytes.NewReader(b))

	b, err = ioutil.ReadAll(w.Body)
	require.NoError(m.T(), err)

	require.Equal(m.T(), http.StatusOK, w.Code)
}

func (m *NetTest) TestUpdateUser() {
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

	b, err := json.Marshal(testUser1)
	require.NoError(m.T(), err)

	w := m.performRequest(http.MethodPut, "/api/v1/user/update", map[string]string{"email": "newEmail", "password": "newPassword"}, bytes.NewReader(b))

	b, err = ioutil.ReadAll(w.Body)
	require.NoError(m.T(), err)

	require.Equal(m.T(), http.StatusOK, w.Code)
}

func (m *NetTest) TestDeleteUser() {
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

	b, err := json.Marshal(testUser1)
	require.NoError(m.T(), err)

	w := m.performRequest(http.MethodDelete, "/api/v1/user/delete", nil, bytes.NewReader(b))

	b, err = ioutil.ReadAll(w.Body)
	require.NoError(m.T(), err)

	require.Equal(m.T(), http.StatusOK, w.Code)
}

func (m *NetTest) createRouter() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v1 := r.Group("/api/v1")

	u := v1.Group("/user")
	u.POST("/create", CreateUser)
	u.GET("/balance", GetUserBalance)
	u.POST("/transfer/:receiver/:amount", TransferFromBalance)
	u.PUT("/update", UpdateUser)
	u.DELETE("/delete", DeleteUser)

	m.router = r
}

func (m *NetTest) performRequest(method, path string, paramKeyVal map[string]string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Printf("Failed to make request %s", err.Error())
	}

	if req != nil {
		q := req.URL.Query()
		if len(paramKeyVal) > 0 {
			for key, val := range paramKeyVal {
				q.Add(key, val)
			}
		}
		req.Header.Add("Content-Type", "application/json")
		req.URL.RawQuery = q.Encode()
	}

	w := httptest.NewRecorder()

	m.router.ServeHTTP(w, req)

	return w
}
