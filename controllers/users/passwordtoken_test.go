package users

import (
	"errors"
	"fmt"
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

func TestGetPasswordResetToken(t *testing.T) {
	tests := []struct {
		Name       string
		idString   string
		AdminUser  auth.TokenData
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{

		{
			"GetRegistrationToken - User Missing from Context",
			"2",
			auth.TokenData{},
			nil,
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetRegistrationToken - non int id param",
			"not-an-int",
			auth.TokenData{UserID: 1},
			nil,
			http.StatusBadRequest,
			`{"error": "Invalid ID"}`,
		},
		{
			"GetPasswordResetToken - SQL Error on InsertRegistrationToken",
			"2",
			auth.TokenData{UserID: 1},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users_tokens").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetPasswordResetToken - Success",
			"2",
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
			r.GET("/api/v1/admin/users/:userId/resettoken", func(c *gin.Context) {
				// Add user to context if it exists
				if test.AdminUser.UserID != 0 {
					c.Set("user", test.AdminUser)
				}

				// Call handler
				GetPasswordResetToken(env)(c)
			})

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/admin/users/%s/resettoken", test.idString), nil)
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
