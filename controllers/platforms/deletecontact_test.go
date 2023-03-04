package platforms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestDeleteContact(t *testing.T) {
	tests := []struct {
		Name             string
		IdString         string
		PlatformIdString string
		MockDbCall       func(sqlmock.Sqlmock)
		StatusCode       int
		Response         string
	}{
		{
			"DeleteContact - non int id",
			"notanint",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{"error":"Invalid ID"}`,
		},
		{
			"DeleteContact - non int platformId",
			"1",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{"error":"Invalid Platform ID"}`,
		},
		{
			"DeleteContact - sql error on DeleteContact",
			"1",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE contacts SET").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"DeleteContact - Valid Request",
			"1",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE contacts SET").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"message": "Contact deleted successfully"}`,
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
			r.DELETE("/api/v1/platforms/:platformId/contacts/:id", DeleteContact(env))

			// Create httptest request
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/platforms/%s/contacts/%s", test.PlatformIdString, test.IdString), nil)
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
