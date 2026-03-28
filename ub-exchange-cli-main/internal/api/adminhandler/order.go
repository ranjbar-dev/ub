package adminhandler

import (
	"exchange-go/internal/api/handler"
	"exchange-go/internal/order"

	"github.com/gin-gonic/gin"
)

func FulFillOrder(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.FulfillOrderParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := handler.HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		u, ok := handler.GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.FulfillOrder(u, p)
		c.JSON(statusCode, resp)
	}
}
