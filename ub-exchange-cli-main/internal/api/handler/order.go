package handler

import (
	"exchange-go/internal/order"
	"exchange-go/internal/userdevice"
	"strings"

	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
)

func CreateOrder(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.CreateOrderParams{}
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

		userAgentHeader := c.GetHeader(HeaderUserAgent)
		device := userdevice.GetDeviceUsingUserAgent(userAgentHeader)

		browserName := ""
		if device == userdevice.DeviceWeb {
			ua := uasurfer.Parse(userAgentHeader)
			browserName = ua.Browser.Name.String()
		}

		uai := order.UserAgentInfo{
			IP:      GetClientIP(c),
			Device:  strings.ToLower(device),
			Browser: browserName,
		}

		p.UserAgentInfo = uai
		resp, statusCode := s.CreateOrder(u, p)
		c.JSON(statusCode, resp)
	}

}

func CancelOrder(s order.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.CancelOrder)
}

func OpenOrders(s order.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetOpenOrders)
}

func OrdersHistory(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.GetOrdersHistoryParams{}
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
		resp, statusCode := s.GetOrdersHistory(u, p, false)
		c.JSON(statusCode, resp)
	}

}

func FullOrdersHistory(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.GetOrdersHistoryParams{}
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

		resp, statusCode := s.GetOrdersHistory(u, p, true)
		c.JSON(statusCode, resp)
	}

}

func TradesHistory(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.GetTradesHistoryParams{}
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

		resp, statusCode := s.GetTradesHistory(u, p, false)
		c.JSON(statusCode, resp)
	}

}

func FullTradesHistory(s order.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := order.GetTradesHistoryParams{}
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

		resp, statusCode := s.GetTradesHistory(u, p, true)
		c.JSON(statusCode, resp)
	}

}

func GetOrderDetail(s order.Service) gin.HandlerFunc {
	return AuthBindQueryAndCall(s.GetOrderDetail)
}
