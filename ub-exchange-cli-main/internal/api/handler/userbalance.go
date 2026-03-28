package handler

import (
	"exchange-go/internal/userbalance"

	"github.com/gin-gonic/gin"
)

func PairBalances(s userbalance.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := userbalance.GetPairBalancesParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.GetPairBalances(u, p)
		c.JSON(statusCode, resp)
	}

}

func AllBalances(s userbalance.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := userbalance.GetAllBalancesParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		u, ok := GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.GetAllBalances(u, p)
		c.JSON(statusCode, resp)
	}

}

func WithdrawAndDeposit(s userbalance.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := userbalance.GetWithdrawDepositParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		u, ok := GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.GetWithdrawDepositData(u, p)
		c.JSON(statusCode, resp)
	}
}

func SetAutoExchange(s userbalance.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.SetAutoExchangeCoin)
}
