package handler

import (
	"exchange-go/internal/country"
	"github.com/gin-gonic/gin"
)

func Countries(s country.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCode := s.GetCountries()
		c.JSON(statusCode, resp)
	}
}
