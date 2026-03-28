package handler

import (
	"exchange-go/internal/configuration"
	"exchange-go/internal/response"
	"github.com/gin-gonic/gin"
)

func GetRecaptchaKey(s configuration.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := configuration.GetRecaptchaKeyParams{}
		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		resp, statusCode := s.GetRecaptchaKey(p)
		c.JSON(statusCode, resp)
	}
}

func GetAppVersion(s configuration.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := configuration.AppVersionParams{}
		err := c.ShouldBindQuery(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := s.GetAppVersion(p)
		c.JSON(statusCode, resp)
	}
}

func ContactUs(s configuration.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := configuration.ContactUsParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := s.ContactUs(p)
		c.JSON(statusCode, resp)
	}
}

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, statusCOde := response.Success("", "OK")
		c.JSON(statusCOde, resp)
	}
}
