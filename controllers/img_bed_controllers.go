package controllers

import (
	"fmt"
	"hugo_backend/config"
	"hugo_backend/utils"
	"image"
	_ "image/jpeg" // 注册JPEG解码器
	_ "image/png"  // 注册PNG解码器
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gen2brain/avif"

	"github.com/gin-gonic/gin"
)

// 定义一个常量作为秘钥（在实际应用中，请从配置文件或环境变量中获取）
// var secretKey string = config.GlobalConfigInstance.Token

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	// 解码图片
	img, format, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	// 检查文件格式
	if format != "jpeg" && format != "png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG formats are supported"})
		return
	}

	// 获取当前时间
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	// 构建目录路径
	uploadPath := filepath.Join("./static/api/i", year, month, day)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 构建 AVIF 文件路径
	randomString, err := utils.GenerateRandomString(6)
	if err != nil {
		log.Println("构建随机字符串失败")
	}
	timestamp := now.UnixNano()
	fileName := randomString + strconv.FormatInt(timestamp, 10)
	avifFileName := fmt.Sprintf("%s.avif", fileName)
	avifFilePath := filepath.Join(uploadPath, avifFileName)

	// 创建 AVIF 文件
	out, err := os.Create(avifFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create avif file"})
		return
	}
	defer out.Close()
	avifOptions := avif.Options{
		Quality:           85, // 降低质量
		QualityAlpha:      85,
		Speed:             10, // 提高编码速度
		ChromaSubsampling: image.YCbCrSubsampleRatio444,
	}
	// 编码并保存图片为 Avif 格式
	err = avif.Encode(out, img, avifOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image to avif"})
		return
	}

	// 返回 avif 文件的 URL
	imageURL := fmt.Sprintf("%s/api/i/%s/%s/%s/%s", config.GlobalConfigInstance.Domain, year, month, day, avifFileName)
	c.JSON(http.StatusOK, gin.H{"result": "success", "code": http.StatusOK, "url": imageURL})
}

func GetImage(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	day := c.Param("day")
	filename := c.Param("filename")

	filePath := filepath.Join("./static/api/i", year, month, day, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"result": "false", "code": http.StatusNotFound, "error": "Image not found"})
		return
	}

	c.File(filePath)
}
