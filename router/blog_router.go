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

func SetSessionCookieMiddleWare(router *gin.Engine) {
	// 使用 cookie 存储会话，并设置密钥
	store := cookie.NewStore([]byte("your-secret-key"))
	router.Use(sessions.Sessions("mysession", store))
}
func SetBlogBackEndTemplateRoutes(templateEmbed embed.FS, router *gin.Engine) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			a, err := json.Marshal(v)
			return string(a), err
		},
	}).ParseFS(templateEmbed, "templates/*.html"))
	router.SetHTMLTemplate(tmpl)
}
func SetBlogBackEndStaticRoutes(staticFiles embed.FS, router *gin.Engine) {
	// 设置静态文件路由
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	// 设置静态文件路由
	router.StaticFS("/api/static", http.FS(staticFS))
}
func SetBlogBackEndRoutes(router *gin.Engine) {
	// Routes for rendering HTML pages
	router.GET("/api/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/api/login", controllers.LoginHandler)
	// API routes
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
