package handler

import (
	"exchange-go/internal/payment"

	"github.com/gin-gonic/gin"
)

func GetPayments(s payment.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetPayments)
}

func GetPaymentDetail(s payment.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.Detail)
}

func PreWithdraw(s payment.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.PreWithdraw)
}

func Withdraw(s payment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := payment.WithdrawParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		u, ok := GetAuthUser(c)
		if !ok {
			return
		}
		ip := GetClientIP(c)
		p.IP = ip
		resp, statusCode := s.Withdraw(u, p)
		c.JSON(statusCode, resp)
	}

}

func Cancel(s payment.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.CancelWithdraw)
}
