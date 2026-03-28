package handler

import (
	"exchange-go/internal/currency"

	"github.com/gin-gonic/gin"
)

func GetCurrencies(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCode := s.GetCurrenciesList()
		c.JSON(statusCode, resp)
	}
}

func GetPairs(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCode := s.GetPairs()
		c.JSON(statusCode, resp)
	}
}

func GetPairsStatistic(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := currency.GetPairsStatisticParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		resp, statusCode := s.GetPairsStatistic(p)
		c.JSON(statusCode, resp)
	}
}

func AddOrRemoveFavoritePair(s currency.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.AddOrRemoveFavoritePair)
}

func GetFavoritePairs(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}
		resp, statusCode := s.GetFavoritePairs(u)
		c.JSON(statusCode, resp)
	}
}


func GetPairRatio(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := currency.GetPairRatioParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		resp, statusCode := s.GetPairRatio(p)
		c.JSON(statusCode, resp)
	}
}

func GetPairsList(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCode := s.GetPairsList()
		c.JSON(statusCode, resp)
	}
}

func GetFees(s currency.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCode := s.GetFees()
		c.JSON(statusCode, resp)
	}
}