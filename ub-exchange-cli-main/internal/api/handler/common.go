package handler

import (
	"exchange-go/internal/api/middleware"
	"exchange-go/internal/user"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BindAndCall handles the common pattern: bind JSON params → call service → return response.
// Use for public endpoints that do not require authentication.
func BindAndCall[P any, R any](fn func(P) (R, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p P
		if err := c.ShouldBindJSON(&p); err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		resp, statusCode := fn(p)
		c.JSON(statusCode, resp)
	}
}

// AuthBindAndCall handles the common pattern: bind JSON params → get auth user → call service → return response.
// Use for authenticated endpoints with a JSON request body.
func AuthBindAndCall[P any, R any](fn func(*user.User, P) (R, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p P
		if err := c.ShouldBindJSON(&p); err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}
		resp, statusCode := fn(u, p)
		c.JSON(statusCode, resp)
	}
}

// AuthCall handles authenticated endpoints with no JSON body (e.g. GET requests with no params).
func AuthCall[R any](fn func(*user.User) (R, int)) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}
		resp, statusCode := fn(u)
		c.JSON(statusCode, resp)
	}
}

const (
	// HeaderUserAgent is the standard User-Agent HTTP header name.
	HeaderUserAgent = "User-Agent"
	// HeaderXForwardedFor is the proxy-forwarded client IP header.
	HeaderXForwardedFor = "x-forwarded-for"
	// HeaderAuthorization is the standard Authorization HTTP header.
	HeaderAuthorization = "Authorization"
)

// GetAuthUser safely extracts the authenticated *user.User from the Gin context.
// Returns the user pointer and true on success. On failure it aborts the request
// with a 401 JSON response and returns nil, false.
func GetAuthUser(c *gin.Context) (*user.User, bool) {
	val, exists := c.Get(middleware.UserKey)
	if !exists || val == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
			Status:  false,
			Message: "authentication required",
		})
		return nil, false
	}
	u, ok := val.(*user.User)
	if !ok || u == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
			Status:  false,
			Message: "invalid user session",
		})
		return nil, false
	}
	return u, true
}

// GetClientIP extracts the real client IP from the X-Forwarded-For header,
// falling back to Gin's ClientIP() if the header is absent or invalid.
func GetClientIP(c *gin.Context) string {
	forwardHeader := c.Request.Header.Get(HeaderXForwardedFor)
	firstAddress := strings.Split(forwardHeader, ",")[0]
	if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
		return firstAddress
	}
	return c.ClientIP()
}
