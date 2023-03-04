package platforms

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

func TestCreateContact(t *testing.T) {
	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"CreateContact - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Invalid ID"}`,
		},
		{
			"CreateContact - Bad json body",
			"1",
			nil,
			http.StatusBadRequest,
			`{badbody}`,
			`{"error": "invalid character 'b' looking for beginning of object key string"}`,
		},
		{
			"CreateContact - sql error on InsertContact",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO contacts").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test","title":"test","email":"test","phone":"","phone2":"","address":"","notes":"","source":"test","privacy":"test"}`,
			`{}`,
		},
		{
			"CreateContact - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO contacts").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"name":"test","title":"test","email":"test","phone":"","phone2":"","address":"","notes":"","source":"test","privacy":"test"}`,
			`{"message":"Contact created successfully"}`,
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
			r.POST("/api/v1/platforms/:platformId/contacts", CreateContact(env))

			// Create httptest request
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/platforms/%s/contacts", test.IdString), strings.NewReader(test.Body))
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
