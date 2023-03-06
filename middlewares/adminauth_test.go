package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/auth"
)

func TestAdminAuthMiddleware(t *testing.T) {
	tests := []struct {
		name string
		user auth.TokenData
		code int
	}{
		{
			"No user token - Unauthorized",
			auth.TokenData{},
			http.StatusUnauthorized,
		},
		{
			"Non admin token - Forbidden",
			auth.TokenData{
				UserID: 1,
				Role:   "test",
			},
			http.StatusForbidden,
		},
		{
			"admin token - OK",
			auth.TokenData{
				UserID: 1,
				Role:   "admin",
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test conext
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Add the query parameters to the Request of the test context
			ctx.Request = &http.Request{}

			if tt.user.Role != "" {
				ctx.Set("user", tt.user)
			}

			// Call middleware on the test context
			AdminAuthMiddleware()(ctx)

			// Make assertions on expected page and pageSize
			require.Equal(t, tt.code, ctx.Writer.Status())
		})
	}
}
