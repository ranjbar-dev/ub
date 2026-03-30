package handler

import (
	"exchange-go/internal/api/middleware"
	"exchange-go/internal/auth"
	"exchange-go/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type centrifugoSubscribeRequest struct {
	Channel string `json:"channel" binding:"required"`
}

func CentrifugoConnectionToken(s auth.CentrifugoTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(middleware.UserKey).(*user.User)
		resp, statusCode := s.GenerateConnectionToken(u.ID)
		c.JSON(statusCode, resp)
	}
}

func CentrifugoSubscriptionToken(s auth.CentrifugoTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(middleware.UserKey).(*user.User)
		p := centrifugoSubscribeRequest{}
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusUnprocessableEntity, NewErrorResponse("invalid request", nil))
			return
		}
		resp, statusCode := s.GenerateSubscriptionToken(u.ID, p.Channel)
		c.JSON(statusCode, resp)
	}
}
