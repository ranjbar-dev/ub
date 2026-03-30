package handler

import (
	"exchange-go/internal/userbalance"

	"github.com/gin-gonic/gin"
)

func PairBalances(s userbalance.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetPairBalances)
}

func AllBalances(s userbalance.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetAllBalances)
}

func WithdrawAndDeposit(s userbalance.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetWithdrawDepositData)
}

func SetAutoExchange(s userbalance.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.SetAutoExchangeCoin)
}
