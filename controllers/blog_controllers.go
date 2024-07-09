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

func HugoUpdateHandler(c *gin.Context) {
	var req models.HugoUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Updating file: %s", req.FilePath)

	if err := os.WriteFile(filepath.Join("content", "posts", req.FilePath), []byte(req.Content), 0644); err != nil {
		log.Printf("Failed to write file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
func HugoDeleteHandler(c *gin.Context) {
	var req models.HugoUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join("content", "posts", req.FilePath)
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func HugoBuildHandler(c *gin.Context) {

	deployCommand := exec.Command("./deploy.sh")
	if err := deployCommand.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetContentHandler(c *gin.Context) {
	filePath := c.Query("filePath")
	content, err := os.ReadFile(filepath.Join("content", "posts", filePath))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

func ListPostsHandler(c *gin.Context) {
	rootDir := "content/posts"
	tree, err := utils.BuildDirectoryTree(rootDir, rootDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{"tree": tree})
}
func HugoNewHandler(c *gin.Context) {
	var req models.HugoNewRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := exec.Command("hugo", "new", "posts/"+req.ContentName+".md")
	if err := command.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filePath := filepath.Join("content", "posts", req.ContentName+".md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileName := req.ContentName + ".md"

	c.JSON(http.StatusOK, gin.H{"content": string(content), "filePath": fileName})
}
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != config.GlobalConfigInstance.UserName || password != config.GlobalConfigInstance.PassWord {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("user", username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
