package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"hugo_backend/config"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

type HugoNewRequest struct {
	ContentName string `json:"contentName"`
}

type HugoUpdateRequest struct {
	FilePath string `json:"filePath"`
	Content  string `json:"content"`
}

type PostContent struct {
	Content string `json:"content"`
}

type Post struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Directory struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Posts    []Post      `json:"posts"`
	Children []Directory `json:"children"`
}

func main() {
	router := gin.Default()
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			a, err := json.Marshal(v)
			return string(a), err
		},
	}).ParseFS(templates, "templates/*.html"))

	router.SetHTMLTemplate(tmpl)
	// 使用 cookie 存储会话，并设置密钥
	store := cookie.NewStore([]byte("your-secret-key"))
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/api", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/api/list")
	})
	// Routes for rendering HTML pages
	router.GET("/api/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/api/login", loginHandler)
	// API routes
	auth := router.Group("/api")
	auth.Use(authRequired)
	{
		auth.GET("/list", listPostsHandler)
		auth.GET("/update", func(c *gin.Context) {
			c.HTML(http.StatusOK, "update.html", nil)
		})
		auth.POST("/hugo/new", hugoNewHandler)
		auth.POST("/hugo/update", hugoUpdateHandler)
		auth.POST("/hugo/build", hugoBuildHandler)
		auth.GET("/getContent", getContentHandler)
		auth.POST("/hugo/delete", hugoDeleteHandler)
	}
	log.Println("Server started at :3000")
	log.Fatal(router.Run(":3000"))
}
func loginHandler(c *gin.Context) {
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

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/api/login")
		c.Abort()
		return
	}

	c.Next()
}
func hugoNewHandler(c *gin.Context) {
	var req HugoNewRequest
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

func hugoUpdateHandler(c *gin.Context) {
	var req HugoUpdateRequest
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
func hugoDeleteHandler(c *gin.Context) {
	var req HugoUpdateRequest
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

func hugoBuildHandler(c *gin.Context) {

	deployCommand := exec.Command("./deploy.sh")
	if err := deployCommand.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func getContentHandler(c *gin.Context) {
	filePath := c.Query("filePath")
	content, err := os.ReadFile(filepath.Join("content", "posts", filePath))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(content)})
}

func listPostsHandler(c *gin.Context) {
	rootDir := "content/posts"
	tree, err := buildDirectoryTree(rootDir, rootDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{"tree": tree})
}

func buildDirectoryTree(root, dir string) (Directory, error) {
	var tree Directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return tree, err
	}

	tree.Name = filepath.Base(dir)
	tree.Path = strings.TrimPrefix(dir, root+"/")
	for _, file := range files {
		if file.IsDir() {
			childTree, err := buildDirectoryTree(root, filepath.Join(dir, file.Name()))
			if err != nil {
				return tree, err
			}
			tree.Children = append(tree.Children, childTree)
		} else if strings.HasSuffix(file.Name(), ".md") {
			post := Post{
				Name: file.Name(),
				Path: strings.TrimPrefix(filepath.Join(dir, file.Name()), root+"/"),
			}
			tree.Posts = append(tree.Posts, post)
		}
	}
	return tree, nil
}
