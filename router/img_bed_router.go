package router

import (
	"hugo_backend/controllers"
	"hugo_backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpImageBedRoute(router *gin.Engine) {
	ImageBedGroup := router.Group("/api")
	ImageBedGroup.Use(middleware.AuthRequired)
	{
		ImageBedGroup.POST("/upload", controllers.UploadImage)
		ImageBedGroup.GET("/i/:year/:month/:day/:filename", controllers.GetImage)
	}
}
