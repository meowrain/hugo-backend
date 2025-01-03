package router

import (
	"embed"
	"encoding/json"
	"html/template"
	"hugo_backend/controllers"
	"hugo_backend/middleware"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 设置会话Cookie中间件
// 为Gin路由器配置会话中间件。
// 它使用基于cookie的会话存储，并设置一个密钥用于会话数据加密。
func SetSessionCookieMiddleWare(router *gin.Engine) {
	// 使用cookie存储会话，并设置密钥
	store := cookie.NewStore([]byte("your-secret-key"))
	router.Use(sessions.Sessions("mysession", store))
}

// 设置博客后端模板路由
// 为Gin路由器配置模板路由。
// 它解析并从嵌入的文件系统加载HTML模板，并将其设置为HTML模板引擎。
func SetBlogBackEndTemplateRoutes(templateEmbed embed.FS, router *gin.Engine) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			a, err := json.Marshal(v)
			return string(a), err
		},
	}).ParseFS(templateEmbed, "templates/*.html"))
	router.SetHTMLTemplate(tmpl)
}

// 设置博客后端静态文件路由
// 为Gin路由器配置静态文件路由。
// 它从嵌入的文件系统中提取静态文件，并将其挂载到指定的路径。
func SetBlogBackEndStaticRoutes(staticFiles embed.FS, router *gin.Engine) {
	// 设置静态文件路由
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	// 设置静态文件路由
	router.StaticFS("/api/static", http.FS(staticFS))
}

// 设置博客后端路由
// 为Gin路由器配置博客后端的路由。
// 它包括用于渲染HTML页面的路由和API路由。
func SetBlogBackEndRoutes(router *gin.Engine) {
	// 渲染HTML页面的路由
	router.GET("/api/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/api/login", controllers.LoginHandler)

	// API路由
	auth := router.Group("/api")
	auth.Use(middleware.AuthRequired)
	{
		auth.GET("/list", controllers.ListPostsHandler)
		auth.GET("/update", func(c *gin.Context) {
			c.HTML(http.StatusOK, "update.html", nil)
		})
		auth.POST("/hugo/new", controllers.HugoNewHandler)
		auth.POST("/hugo/update", controllers.HugoUpdateHandler)
		auth.POST("/hugo/build", controllers.HugoBuildHandler)
		auth.GET("/getContent", controllers.GetContentHandler)
		auth.POST("/hugo/delete", controllers.HugoDeleteHandler)
	}
}
