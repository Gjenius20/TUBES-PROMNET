package middleware

import (
	"net/http"
	"strings"

	"autograder/internal/models"
	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserID = "userID"
	ctxRole   = "role"
)

// Auth memverifikasi access token JWT dari header "Authorization: Bearer <token>".
func Auth(tokens *utils.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			utils.Error(c, http.StatusUnauthorized, "missing or malformed authorization header")
			c.Abort()
			return
		}
		claims, err := tokens.ParseAccessToken(strings.TrimSpace(parts[1]))
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "invalid or expired access token")
			c.Abort()
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// RequireRole membatasi route hanya untuk role tertentu (harus dipasang setelah Auth).
func RequireRole(role models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		if CurrentRole(c) != role {
			utils.Error(c, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) uint {
	v, _ := c.Get(ctxUserID)
	id, _ := v.(uint)
	return id
}

func CurrentRole(c *gin.Context) models.Role {
	v, _ := c.Get(ctxRole)
	role, _ := v.(string)
	return models.Role(role)
}
