package projects

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

func TestEditProject(t *testing.T) {
	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"EditProject - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Invalid ID"}`,
		},
		{
			"EditProject - Bad json body",
			"1",
			nil,
			http.StatusBadRequest,
			`{badbody}`,
			`{"error": "invalid character 'b' looking for beginning of object key string"}`,
		},
		{
			"EditProject - sql error on EditProject transaction",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE projects SET").WillReturnError(errors.New("test"))
				mock.ExpectRollback()
			},
			http.StatusInternalServerError,
			`{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test"}`,
			`{}`,
		},
		{
			"EditProject - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE projects SET").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM projects_tags").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO projects_tags").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM platforms_projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO platforms_projects").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			http.StatusOK,
			`{"title":"test","description":"test","link":"test","date":{"Time":"2023-03-05T00:00:00Z","Valid":true},"body":"test","tags":[{"id":1,"tag":"test"}],"platforms":[{"id":1,"platform":"test"}]}`,
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
			r.PUT("/api/v1/projects/:projectId", EditProject(env))

			// Create httptest request
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/projects/%s", test.IdString), strings.NewReader(test.Body))
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
