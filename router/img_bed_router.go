package router

import (
	"hugo_backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpImageBedRoute(router *gin.Engine) {
	ImageBedGroup := router.Group("/api")
	{
		ImageBedGroup.GET("/i/:year/:month/:day/:filename", controllers.GetImage)
		ImageBedGroup.POST("/upload", controllers.UploadImage)
	}
}
