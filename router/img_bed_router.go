package router

import (
	"github.com/gin-gonic/gin"
	"hugo_backend/controllers"
)

// SetUpImageBedRoute 设置图片床服务的路由。
// 它为图片床路由创建一个新的组，并添加必要的中间件和处理程序。
// 路由前缀为 "/api"。
func SetUpImageBedRoute(router *gin.Engine) {
	// 为图片床路由创建一个新的组
	ImageBedGroup := router.Group("/api")

	// 添加一个 GET 路由，通过文件名和日期获取图片
	ImageBedGroup.GET("/i/:year/:month/:day/:filename", controllers.GetImage)

	// 为该组中的所有路由使用 AuthRequired 中间件
	//ImageBedGroup.Use(middleware.AuthRequired)

	// 添加一个 POST 路由，用于上传图片
	ImageBedGroup.POST("/upload", controllers.UploadImage)
}
