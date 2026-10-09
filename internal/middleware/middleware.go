package middleware

import (
	"net/http"
	"strings"

	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ any) {
		utils.Fail(c, http.StatusInternalServerError, "internal server error")
		c.Abort()
	})
}

// Auth verifies the access token on every route it wraps.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		tok := strings.TrimPrefix(h, "Bearer ")
		if h == "" || tok == h || tok == "" {
			utils.Fail(c, http.StatusUnauthorized, "missing or malformed Authorization header")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(tok, secret)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("uid", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r, _ := c.Get("role"); r != role {
			utils.Fail(c, http.StatusForbidden, "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
