package controllers

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
	"github.com/webstradev/rsdb-backend/utils"
)

func TestEditPlatform(t *testing.T) {
	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"EditPlatform - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Invalid ID"}`,
		},
		{
			"EditPlatform - missing required fields",
			"1",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Key: 'editPlatformInput.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'editPlatformInput.Country' Error:Field validation for 'Country' failed on the 'required' tag\nKey: 'editPlatformInput.Privacy' Error:Field validation for 'Privacy' failed on the 'required' tag\nKey: 'editPlatformInput.Categories' Error:Field validation for 'Categories' failed on the 'required' tag"}`,
		},
		{
			"EditPlatform - sql error on EditPlatform",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE platforms SET").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[]}`,
			`{}`,
		},
		{
			"EditPlatform - sql error on UpdatePlatformCategories(DELETE)",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE platforms SET").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM platforms_categories").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[]}`,
			`{}`,
		},
		{
			"EditPlatform - sql error on UpdatePlatformCategories(INSERT)",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE platforms SET").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM platforms_categories").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO platforms_categories").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test", "country":"test", "privacy":"Private", "categories":[1]}`,
			`{}`,
		},
		{
			"EditPlatform - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE platforms SET").WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("DELETE FROM platforms_categories").WillReturnResult(sqlmock.NewResult(1, 1))

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
			r.PUT("/api/v1/platforms/:platformId", EditPlatform(env))

			// Create httptest request
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/platforms/%s", test.IdString), strings.NewReader(test.Body))
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
