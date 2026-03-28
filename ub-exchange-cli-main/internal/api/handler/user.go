package handler

import (
	"exchange-go/internal/user"
	"exchange-go/internal/userdevice"
	"strings"

	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
)

func SetUserProfile(s user.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.SetUserProfile)
}

func GetUserProfile(s user.Service) gin.HandlerFunc {
	return AuthCall(s.GetProfile)
}

func GetUserData(s user.Service) gin.HandlerFunc {
	return AuthCall(s.GetUserData)
}

func Get2FaBarcode(s user.Service) gin.HandlerFunc {
	return AuthCall(s.Get2FaBarcode)
}

func Enable2Fa(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := user.Enable2FaParams{}
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
		p.IP = GetClientIP(c)
		p.UserAgent = userAgentHeader
		resp, statusCode := s.Enable2Fa(u, p)
		c.JSON(statusCode, resp)
	}
}

func Disable2Fa(s user.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.Disable2Fa)
}

func ChangePassword(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := user.ChangePasswordParams{}
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
		ip := GetClientIP(c)
		uai := user.RequestUserAgentInfo{
			IP:      ip,
			Device:  strings.ToLower(device),
			Browser: browserName,
		}
		p.IP = ip
		p.UserAgent = userAgentHeader

		p.UserAgentInfo = uai

		resp, statusCode := s.ChangePassword(u, p)
		c.JSON(statusCode, resp)
	}
}

func SendSms(s user.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.SendSms)
}

func EnableSms(s user.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.EnableSms)
}

func DisableSms(s user.Service) gin.HandlerFunc {
	return AuthBindAndCall(s.DisableSms)
}

func SendVerificationEmail(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.SendVerificationEmail(u)
		c.JSON(statusCode, resp)
	}
}
