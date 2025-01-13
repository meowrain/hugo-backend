package main

import (
	"embed"
	"hugo_backend/router"
	"log"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var staticFiles embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 设置模板
	router.SetBlogBackEndTemplateRoutes(templates, r)
	// 静态文件路由
	router.SetBlogBackEndStaticRoutes(staticFiles, r)
	// 使用 cookie 存储会话，并设置密钥
	router.SetSessionCookieMiddleWare(r)
	//设置路由
	router.SetBlogBackEndRoutes(r) //设置博客路由
	router.SetUpImageBedRoute(r)   //设置图床路由
	log.Println("Server started at :3000")
	log.Fatal(r.Run(":3000"))
}
