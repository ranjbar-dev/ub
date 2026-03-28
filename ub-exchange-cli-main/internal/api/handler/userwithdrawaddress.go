package handler

import (
	"exchange-go/internal/userwithdrawaddress"

	"github.com/gin-gonic/gin"
)

func GetWithdrawAddresses(s userwithdrawaddress.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := userwithdrawaddress.GetWithdrawAddressesParams{}
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

		resp, statusCode := s.GetWithdrawAddresses(u, p)
		c.JSON(statusCode, resp)
	}

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
	return func(c *gin.Context) {
		p := userwithdrawaddress.GetFormerAddressesParams{}
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

		resp, statusCode := s.GetFormerAddresses(u, p)
		c.JSON(statusCode, resp)
	}

}
