package middleware

import (
	"github.com/gin-gonic/gin"
)

const NonRequiredAuthKey = "nonRequiredAuth"

func NonRequiredAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(NonRequiredAuthKey, true)
		c.Next()
		return
	}
}
