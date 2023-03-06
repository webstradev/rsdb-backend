package projects

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetProject(t *testing.T) {
	// This timestamp is to mock date values returned by the database
	timestamp, err := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}

	sqlTimestamp := sql.NullTime{
		Time:  timestamp,
		Valid: true,
	}

	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetProject - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
		},
		{
			"GetProject - sql error on GetProject",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT p.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetProject - project not found",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT p.(.+)").WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			http.StatusNotFound,
			`{}`,
		},
		{
			"GetProject - sql error on GetProjectTags",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sqlTimestamp, "test")
				mock.ExpectQuery("SELECT p.(.+)").WithArgs(1).WillReturnRows(rows)

				mock.ExpectQuery("SELECT pt.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetProject - sql error on GetProjectPlatforms",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sqlTimestamp, "test")
				mock.ExpectQuery("SELECT p.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"project_id", "tag_id", "tag"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT pt.(.+)").WithArgs(1).WillReturnRows(rows)

				mock.ExpectQuery("SELECT pp.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetProject - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sqlTimestamp, "test")
				mock.ExpectQuery("SELECT p.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"project_id", "tag_id", "tag"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT pt.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"project_id", "platform_id", "platform_name"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT pp.(.+)").WithArgs(1).WillReturnRows(rows)
			},
			http.StatusOK,
			`{"id":1,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","tags":[{"id":1,"tag":"test"}],"platforms":[{"id":1,"platform":"test"}]}`,
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
			r.GET("/api/projects/:projectId", GetProject(env))

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/projects/%s", test.IdString), nil)
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
