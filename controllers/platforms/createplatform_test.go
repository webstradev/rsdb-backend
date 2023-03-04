package platforms

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestCreatePlatform(t *testing.T) {
	tests := []struct {
		Name       string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"CreatePlatform - missing required fields",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Key: 'createPlatformInput.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'createPlatformInput.Country' Error:Field validation for 'Country' failed on the 'required' tag\nKey: 'createPlatformInput.Privacy' Error:Field validation for 'Privacy' failed on the 'required' tag\nKey: 'createPlatformInput.Categories' Error:Field validation for 'Categories' failed on the 'required' tag"}`,
		},
		{
			"CreatePlatform - sql error on CreatePlatform",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO platforms").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[]}`,
			`{}`,
		},
		{
			"CreatePlatform - sql error on UpdatePlatformCategories(INSERT)",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO platforms").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO platforms_categories").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[1]}`,
			`{}`,
		},
		{
			"CreatePlatform - Valid Request",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO platforms").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO platforms_categories").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[1]}`,
			`{}`,
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

			// Register handler
			r.POST("/api/v1/platforms", CreatePlatform(env))

			// Create httptest request
			req, _ := http.NewRequest("POST", "/api/v1/platforms", strings.NewReader(test.Body))
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
