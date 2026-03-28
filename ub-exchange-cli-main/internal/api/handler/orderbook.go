package handler

import (
	"exchange-go/internal/orderbook"

	"github.com/gin-gonic/gin"
)

func OrderBook(s orderbook.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := orderbook.GetOrderBookParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := s.GetOrderBook(p)
		c.JSON(statusCode, resp)
	}
}

func TradeBook(s orderbook.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := orderbook.GetTradeBookParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := s.GetTradeBook(p)
		c.JSON(statusCode, resp)
	}

}
