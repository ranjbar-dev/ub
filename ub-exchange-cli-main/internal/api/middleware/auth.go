package middleware

import (
	"exchange-go/internal/auth"
	"exchange-go/internal/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	MessageJwtNotFound          = "JWT Token not found"
	MessageJwtInvalidToken      = "invalid token"
	MessageJwtInvalidCredential = "invalid credential"
	UserKey                     = "user"
)

type authHeader struct {
	Token string `header:"Authorization"`
}


func AuthMiddleware(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			if shouldAbort(c) {
				abort(c, MessageJwtNotFound)
			}
			next(c, nil)
			return
		}

		parts := strings.Split(h.Token, " ")
		if len(parts) < 2 {
			if shouldAbort(c) {
				abort(c, MessageJwtInvalidToken)
			}
			next(c, nil)
			return
		}

		token := parts[1]
		loggedInUser, err := s.GetUser(token)
		if err != nil {
			if shouldAbort(c) {
				abort(c, MessageJwtInvalidToken)
			}
			next(c, nil)
			return

		}

		if loggedInUser != nil {
			next(c, loggedInUser)
			return
		}

		if shouldAbort(c) {
			abort(c, MessageJwtInvalidCredential)
		}
		next(c, loggedInUser)
		return
	}
}

func abort(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  false,
		"message": message,
	})
}

func shouldAbort(c *gin.Context) bool {
	//for nonRequiredAuthKey it means the user can be nil so we do not abort
	// only abort when authentication is required
	return !c.GetBool(NonRequiredAuthKey)
}

func next(c *gin.Context, u *user.User) {
	c.Set(UserKey, u)
	c.Next()
}
