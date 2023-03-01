package controllers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetCounts(t *testing.T) {
	tests := []struct {
		Name       string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetCounts - sql error - CountPlatforms",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetCounts - sql error - CountArticles",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms WHERE deleted_at IS NULL").WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles WHERE deleted_at IS NULL").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetCounts - sql error - CountProjects",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms WHERE deleted_at IS NULL").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(50)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles WHERE deleted_at IS NULL").WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM projects WHERE deleted_at IS NULL").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetCounts - sql error - CountContacts",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms WHERE deleted_at IS NULL").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(50)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles WHERE deleted_at IS NULL").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(200)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM projects WHERE deleted_at IS NULL").WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM contacts WHERE deleted_at IS NULL").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetCounts - successfull",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(50)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(200)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM projects").WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(150)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM contacts WHERE deleted_at IS NULL").WillReturnRows(rows)
			},
			http.StatusOK,
			`{"platforms":100,"contacts":150,"articles":50,"projects":200}`,
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
			r.GET("/api/v1/counts", GetCounts(env))

			// Create httptest request
			req, _ := http.NewRequest("GET", "/api/v1/counts", nil)
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

			// // Check for any remaining expectations
			// // we make sure that all expectations were met
			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}
