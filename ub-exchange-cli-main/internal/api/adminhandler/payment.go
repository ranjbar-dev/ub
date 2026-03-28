package adminhandler

import (
	"exchange-go/internal/api/handler"
	"exchange-go/internal/payment"

	"github.com/gin-gonic/gin"
)

func Callback(s payment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := payment.WalletCallBackParams{}
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

		resp, statusCode := s.HandleWalletCallBack(u, p)
		c.JSON(statusCode, resp)
	}
}

func UpdateWithdraw(s payment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := payment.UpdateWithdrawParams{}
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

		resp, statusCode := s.UpdateWithdraw(u, p)
		c.JSON(statusCode, resp)
	}
}

func UpdateDeposit(s payment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := payment.UpdateDepositParams{}
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

		resp, statusCode := s.UpdateDeposit(u, p)
		c.JSON(statusCode, resp)
	}
}
