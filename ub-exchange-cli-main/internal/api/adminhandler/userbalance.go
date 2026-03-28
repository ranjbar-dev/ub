package adminhandler

import (
	"exchange-go/internal/api/handler"
	"exchange-go/internal/userbalance"

	"github.com/gin-gonic/gin"
)

func UpdateUserBalance(s userbalance.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := userbalance.UpdateUserBalanceFromAdminParams{}
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

		resp, statusCode := s.UpdateUserBalanceFromAdmin(u, p)
		c.JSON(statusCode, resp)
	}
}
