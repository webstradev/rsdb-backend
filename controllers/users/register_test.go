package users

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/mocks"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		Name       string
		Token      string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"Register - missing token",
			"",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Missing token"}`,
		},
		{
			"Register - sql error on ValidateToken",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
			`{}`,
		},
		{
			"Register - invalid token",
			"invalidtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
			},
			http.StatusBadRequest,
			`{}`,
			`{"error": "invalid or expired token"}`,
		},
		{
			"Register - missing email",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
			},
			http.StatusBadRequest,
			`{"password": "test"}`,
			`{"error": "please provide a valid email and a password"}`,
		},
		{
			"Register - missing password",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
			},
			http.StatusBadRequest,
			`{"email": "test"}`,
			`{"error": "please provide a valid email and a password"}`,
		},
		{
			"Register - sql error on IsUsernameAvailable",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
		{
			"Register - username already in use",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
			},
			http.StatusBadRequest,
			`{"email": "test","password":"test"}`,
			`{"error": "an account with this emailadress is already in use"}`,
		},
		{
			"Register - sql error on ConsumeToken",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec("UPDATE users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
		{
			"Register - error on ConsumeToken",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec("UPDATE users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
		{
			"Register - error on CreatePasswordHash",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec("UPDATE users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"error"}`,
			`{}`,
		},
		{
			"Register - error on InsertUser",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec("UPDATE users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO users").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
		{
			"Register - Valid Request",
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectQuery("SELECT COUNT(.+) FROM users").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec("UPDATE users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusAccepted,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Initilize test router, environemnt and mock database
			r, mockDb, mockSql, env, err := utils.SetupTestEnvironment(test.MockDbCall)
			// Close the mock database at the end of the test
			defer mockDb.Close()

			env.AuthService = mocks.NewMockAuthService()

			// Check for errors during setup
			require.NoError(t, err)

			// Register handler
			r.POST("/users/register", Register(env))

			// Create httptest request
			req, _ := http.NewRequest("POST", fmt.Sprintf("/users/register?token=%s", test.Token), strings.NewReader(test.Body))
			w := httptest.NewRecorder()

			// Mock request
			r.ServeHTTP(w, req)

			// Read response data
			responseData, _ := io.ReadAll(w.Body)

			// Check response status
			require.Equal(t, test.StatusCode, w.Code)

			// Handle empty responses
			response := string(responseData)
			if response == "" {
				response = "{}"
			}

			// Check response body
			require.JSONEq(t, test.Response, response)

			// Check for any remaining expectations
			// we make sure that all expectations were met
			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}
