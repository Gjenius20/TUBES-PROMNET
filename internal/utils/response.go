package utils

import "github.com/gin-gonic/gin"

// Success menulis response sukses dengan format standar API.
func Success(c *gin.Context, status int, message string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(status, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error menulis response error dengan format standar API.
func Error(c *gin.Context, status int, message string, errs ...string) {
	if errs == nil {
		errs = []string{}
	}
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"errors":  errs,
	})
}
