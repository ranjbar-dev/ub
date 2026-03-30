package handler

import (
	"exchange-go/internal/userwithdrawaddress"

	"github.com/gin-gonic/gin"
)

func GetWithdrawAddresses(s userwithdrawaddress.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetWithdrawAddresses)
}

func NewWithdrawAddress(s userwithdrawaddress.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.CreateNewAddress)
}

func AddToFavorites(s userwithdrawaddress.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.AddToFavorites)
}

func Delete(s userwithdrawaddress.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.Delete)
}

func GetFormerAddresses(s userwithdrawaddress.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetFormerAddresses)
}
