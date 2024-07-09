package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/api/login")
		c.Abort()
		return
	}

	c.Next()
}
