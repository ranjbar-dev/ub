package handler

import (
	"exchange-go/internal/auth"

	"github.com/gin-gonic/gin"
)

func Login(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.LoginParams{}
		err := c.BindJSON(&p)

		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		p.IP = GetClientIP(c)

		resp, statusCode := s.Login(p)
		c.JSON(statusCode, resp)
	}
}

func Register(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.RegisterParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		p.IP = GetClientIP(c)

		resp, statusCode := s.Register(p)
		c.JSON(statusCode, resp)
	}
}

func ForgotPassword(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.ForgotPasswordParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		p.IP = GetClientIP(c)

		resp, statusCode := s.ForgotPassword(p)
		c.JSON(statusCode, resp)
	}
}

func ForgotPasswordUpdate(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.ForgotPasswordUpdateParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		p.IP = GetClientIP(c)

		resp, statusCode := s.ForgotPasswordUpdate(p)
		c.JSON(statusCode, resp)
	}
}

func VerifyEmail(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.VerifyEmailParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		resp, statusCode := s.VerifyEmail(p)
		c.JSON(statusCode, resp)
	}
}

func GetTokenByRefreshToken(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.GetTokenByRefreshTokenParams{}
		err := c.ShouldBindJSON(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		userAgent := c.GetHeader(HeaderUserAgent)
		p.UserAgent = userAgent
		p.IP = GetClientIP(c)

		resp, statusCode := s.GetTokenByRefreshToken(p)
		c.JSON(statusCode, resp)
	}
}
