package handler

import (
	"exchange-go/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func MqttLogin(s auth.MqttAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//everyone can login so we do not check it right now
		p := auth.MqttLoginParams{}
		resp, statusCode := s.Login(p)
		c.JSON(statusCode, resp)
	}

}
func MqttACL(s auth.MqttAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.MqttACLParams{}
		err := c.ShouldBindWith(&p, binding.FormPost)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		resp, statusCode := s.ACL(p)
		c.JSON(statusCode, resp)
	}

}

func MqttSuperUser(s auth.MqttAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := auth.MqttSuperUserParams{}
		err := c.ShouldBindWith(&p, binding.FormPost)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}

		resp, statusCode := s.SuperUser(p)
		c.JSON(statusCode, resp)
	}

}
