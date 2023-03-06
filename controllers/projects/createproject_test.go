package projects

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

func TestCreateProject(t *testing.T) {
	tests := []struct {
		Name       string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"CreateProject - Bad json body",
			nil,
			http.StatusBadRequest,
			`{badbody}`,
			`{"error": "invalid character 'b' looking for beginning of object key string"}`,
		},
		{
			"CreateProject - sql error on InsertProject",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO projects").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"project":{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test"},"linkedPlatforms":[1,2,3],"tags":[1,2]}`,
			`{}`,
		},
		{
			"CreateProject - sql error on InsertProjectPlatforms",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO platforms_projects").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"project":{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test"},"linkedPlatforms":[1,2,3],"tags":[1,2]}`,
			`{}`,
		},
		{
			"CreateProject - sql error on InsertArticleTags",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO projects_tags").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"project":{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test"},"linkedPlatforms":[],"tags":[1,2]}`,
			`{}`,
		},
		{
			"CreateProject - Valid Request",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO platforms_projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO projects_tags").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"project":{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test"},"linkedPlatforms":[1,2,3],"tags":[1,2]}`,
			`{"message":"Project created successfully"}`,
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
			r.POST("/api/v1/projects", CreateProject(env))

			// Create httptest request
			req, _ := http.NewRequest("POST", "/api/v1/projects", strings.NewReader(test.Body))
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
