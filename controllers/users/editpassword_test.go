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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/mocks"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestEditPassword(t *testing.T) {
	tests := []struct {
		Name       string
		User       auth.TokenData
		Token      string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"EditPassword - missing user from context",
			auth.TokenData{},
			"",
			nil,
			http.StatusInternalServerError,
			`{}`,
			`{}`,
		},
		{
			"EditPassword - missing token",
			auth.TokenData{UserID: 1},
			"",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Missing token"}`,
		},
		{
			"EditPassword - sql error on ValidateToken",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
			`{}`,
		},
		{
			"EditPassword - invalid token",
			auth.TokenData{UserID: 1},
			"invalidtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
			},
			http.StatusBadRequest,
			`{}`,
			`{"error": "invalid or expired token"}`,
		},
		{
			"EditPassword - missing password",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
			},
			http.StatusBadRequest,
			`{}`,
			`{"error": "please provide a new password"}`,
		},
		{
			"EditPassword - sql error on ConsumeToken",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectExec("UPDATE users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test","password":"test"}`,
			`{}`,
		},
		{
			"EditPassword - error on CreatePasswordHash",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
			},
			http.StatusInternalServerError,
			`{"password":"error"}`,
			`{}`,
		},
		{
			"EditPassword - error on UpdateUser",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectExec("UPDATE users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE users").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"password":"test"}`,
			`{}`,
		},
		{
			"EditPassword - Valid Request",
			auth.TokenData{UserID: 1},
			"validtoken",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) FROM users_tokens").WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectExec("UPDATE users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
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

			// EditPassword handler
			r.PUT("/users/password", func(c *gin.Context) {
				// Add user to context if it exists
				if test.User.UserID != 0 {
					c.Set("user", test.User)
				}

				// Call handler
				EditPassword(env)(c)
			})

			// Create httptest request
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/password?token=%s", test.Token), strings.NewReader(test.Body))
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
