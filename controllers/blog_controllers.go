package controllers

import (
	"hugo_backend/config"
	"hugo_backend/models"
	"hugo_backend/utils"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HugoUpdateHandler 处理更新Hugo文章的请求
func HugoUpdateHandler(c *gin.Context) {
	var req models.HugoUpdateRequest
	// 绑定JSON请求到req结构体
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Updating file: %s", req.FilePath)

	// 将更新的内容写入文件
	if err := os.WriteFile(filepath.Join("content", "posts", req.FilePath), []byte(req.Content), 0644); err != nil {
		log.Printf("Failed to write file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功状态码
	c.Status(http.StatusNoContent)
}

// HugoDeleteHandler 处理删除Hugo文章的请求
func HugoDeleteHandler(c *gin.Context) {
	var req models.HugoUpdateRequest
	// 绑定JSON请求到req结构体
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join("content", "posts", req.FilePath)
	// 删除指定文件
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功状态码
	c.Status(http.StatusNoContent)
}

// HugoBuildHandler 处理Hugo项目构建的请求
func HugoBuildHandler(c *gin.Context) {
	// 执行deploy.sh脚本
	deployCommand := exec.Command("./deploy.sh")
	if err := deployCommand.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功状态码
	c.Status(http.StatusNoContent)
}

// GetContentHandler 处理获取文章内容的请求
func GetContentHandler(c *gin.Context) {
	filePath := c.Query("filePath")
	// 读取指定文件的内容
	content, err := os.ReadFile(filepath.Join("content", "posts", filePath))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回文章内容
	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

// ListPostsHandler 处理列出所有文章的请求
func ListPostsHandler(c *gin.Context) {
	rootDir := "content/posts"
	// 构建目录树
	tree, err := utils.BuildDirectoryTree(rootDir, rootDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回目录树的HTML页面
	c.HTML(http.StatusOK, "list.html", gin.H{"tree": tree})
}

// HugoNewHandler 处理创建新文章的请求
func HugoNewHandler(c *gin.Context) {
	var req models.HugoNewRequest
	// 绑定JSON请求到req结构体
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用Hugo命令创建新文章
	command := exec.Command("hugo", "new", "posts/"+req.ContentName+".md")
	if err := command.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join("content", "posts", req.ContentName+".md")
	// 读取新创建文章的内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileName := req.ContentName + ".md"

	// 返回新文章的内容和文件名
	c.JSON(http.StatusOK, gin.H{"content": string(content), "filePath": fileName})
}

// LoginHandler 处理用户登录请求
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名和密码
	if username != config.GlobalConfigInstance.UserName || password != config.GlobalConfigInstance.PassWord {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	// 设置用户会话
	session.Set("user", username)
	session.Save()

	// 返回登录成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
