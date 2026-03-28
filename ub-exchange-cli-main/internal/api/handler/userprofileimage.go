package handler

import (
	"exchange-go/internal/user"

	"github.com/gin-gonic/gin"
)

func MultipleUpload(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := user.UploadImagesParams{}
		err := c.ShouldBind(&p)
		if err != nil {
			errorResponse, statusCode := HandleValidationError(err)
			c.AbortWithStatusJSON(statusCode, errorResponse)
			return
		}
		frontImage, frontHeader, err := c.Request.FormFile("front_image")
		backImage, backHeader, err := c.Request.FormFile("back_image")
		p.FrontImage = frontImage
		p.FrontHeader = frontHeader
		p.BackImage = backImage
		p.BackHeader = backHeader
		u, ok := GetAuthUser(c)
		if !ok {
			return
		}

		resp, statusCode := s.UploadImages(u, p)
		c.JSON(statusCode, resp)
	}
}

func DeleteProfileImage(s user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := user.DeleteImageParams{}
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

		resp, statusCode := s.DeleteImage(u, p)
		c.JSON(statusCode, resp)
	}
}
