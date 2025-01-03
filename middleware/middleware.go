package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired 是一个中间件函数，用于检查用户是否已登录。
// 如果用户未登录，则会重定向到登录页面。
func AuthRequired(c *gin.Context) {
	// 获取当前会话
	session := sessions.Default(c)
	// 从会话中获取用户信息
	user := session.Get("user")

	// 如果用户信息为空，说明用户未登录
	if user == nil {
		// 重定向到登录页面
		c.Redirect(http.StatusTemporaryRedirect, "/api/login")
		// 终止当前请求的处理
		c.Abort()
		return
	}

	// 如果用户已登录，则继续处理请求
	c.Next()
}
