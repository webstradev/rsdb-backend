package articles

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

func TestGetArticle(t *testing.T) {
	// This timestamp is to mock date values returned by the database
	timestamp, err := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetArticle - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
		},
		{
			"GetArticle - sql error on GetArticle",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT a.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetArticle - article not found",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT a.(.+)").WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			http.StatusNotFound,
			`{}`,
		},
		{
			"GetArticle - sql error on GetArticleTags",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+)").WithArgs(1).WillReturnRows(rows)

				mock.ExpectQuery("SELECT at.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetArticle - sql error on GetArticlePlatforms",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"article_id", "tag_id", "tag"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT at.(.+)").WithArgs(1).WillReturnRows(rows)

				mock.ExpectQuery("SELECT pa.(.+)").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetArticle - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"article_id", "tag_id", "tag"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT at.(.+)").WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"article_id", "platform_id", "platform_name"}).
					AddRow(1, 1, "test")
				mock.ExpectQuery("SELECT pa.(.+)").WithArgs(1).WillReturnRows(rows)
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
			r.GET("/api/articles/:articleId", GetArticle(env))

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/articles/%s", test.IdString), nil)
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
