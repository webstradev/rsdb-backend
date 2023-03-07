package users

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/mocks"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetRegistrationToken(t *testing.T) {
	tests := []struct {
		Name       string
		User       auth.TokenData
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{

		{
			"GetRegistrationToken - User Missing from Context",
			auth.TokenData{},
			nil,
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetRegistrationToken - SQL Error on InsertRegistrationToken",
			auth.TokenData{UserID: 1},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetRegistrationToken - Success",
			auth.TokenData{UserID: 1},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users_tokens").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"token": "mock-uuid"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Initilize test router, environemnt and mock database
			r, mockDb, mockSql, env, err := utils.SetupTestEnvironment(test.MockDbCall)
			// Close the mock database at the end of the test
			defer mockDb.Close()

			// Check for errors during setup
			require.NoError(t, err)

			// Add mock uuid service to environment
			env.UUID = mocks.NewMockUUIDService()

			// Register handler
			r.GET("/api/v1/admin/users/token", func(c *gin.Context) {
				// Add user to context if it exists
				if test.User.UserID != 0 {
					c.Set("user", test.User)
				}

				// Call handler
				GetRegistrationToken(env)(c)
			})

			// Create httptest request
			req, _ := http.NewRequest("GET", "/api/v1/admin/users/token", nil)
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
