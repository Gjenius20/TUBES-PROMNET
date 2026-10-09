package middleware

import (
	"log"
	"net/http"

	"autograder/internal/utils"

	"github.com/gin-gonic/gin"
)

// Recovery menangkap panic dan membalas dengan JSON error standar.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("panic recovered: %v", recovered)
		utils.Error(c, http.StatusInternalServerError, "internal server error")
		c.Abort()
	})
}
