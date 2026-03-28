package middleware

import (
	"exchange-go/internal/auth"
	"github.com/gin-gonic/gin"
	"strings"
)

func AdminAuthMiddleware(s auth.Service) gin.HandlerFunc {
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
		loggedInUser, err := s.GetAdminUser(token)
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
