package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/mocks"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestJWTAuthMiddleware(t *testing.T) {

	tests := []struct {
		name              string
		authHeader        string
		want              int
		expectedTokenData *auth.TokenData
	}{
		{
			"No auth header",
			"",
			401,
			nil,
		},
		{
			"Invalid auth header",
			"Bearer invalidtoken",
			401,
			nil,
		},
		{
			"user auth header",
			"Bearer usertoken",
			200,
			&auth.TokenData{UserID: 2, Role: "user"},
		},
		{
			"admin auth header",
			"Bearer admintoken",
			200,
			&auth.TokenData{UserID: 1, Role: "admin"},
		},
	}

	env := utils.Environment{
		JWT: mocks.CreateMockJWTService(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test context
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Create a mock request
			c.Request, _ = http.NewRequest("GET", "/", nil)

			// Use the authorzation header from this test
			c.Request.Header.Set("Authorization", tt.authHeader)

			// Call the middleware on the test context
			JWTAuthMiddleware(&env)(c)

			require.Equal(t, tt.want, c.Writer.Status())

			if tt.expectedTokenData != nil {
				require.Equal(t, *tt.expectedTokenData, c.MustGet("user"))
			}
		})
	}
}
