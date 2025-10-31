package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件大小（限制 5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// 检查文件扩展名（第一层：扩展名白名单）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// 二次 MIME 嗅探：打开文件并读取前 512 字节检测真实 MIME 类型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// 读取文件前 512 字节用于 MIME 检测
	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil && n == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	// 检测真实的 MIME 类型
	detectedMIME := http.DetectContentType(buf[:n])
	if mt, _, err := mime.ParseMediaType(detectedMIME); err == nil {
		detectedMIME = mt
	}

	// 允许的 MIME 类型映射
	allowedMIMEs := map[string]bool{
		"image/jpeg":               true,
		"image/jpg":                true,
		"image/png":                true,
		"image/gif":                true,
		"image/webp":               true,
		"application/octet-stream": true, // 某些图片可能被识别为 octet-stream
	}

	// 验证检测到的 MIME 类型是否在允许列表中
	if !allowedMIMEs[detectedMIME] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid file content. Detected MIME type: %s. Allowed: image/jpeg, image/png, image/gif, image/webp", detectedMIME),
		})
		return
	}

	// 验证扩展名与 MIME 类型的一致性
	extMIMEMap := map[string][]string{
		".jpg":  {"image/jpeg", "image/jpg"},
		".jpeg": {"image/jpeg", "image/jpg"},
		".png":  {"image/png"},
		".gif":  {"image/gif"},
		".webp": {"image/webp"},
	}

	expectedMIMEs := extMIMEMap[ext]
	mimeMatches := false
	for _, expectedMIME := range expectedMIMEs {
		if detectedMIME == expectedMIME || detectedMIME == "application/octet-stream" {
			mimeMatches = true
			break
		}
	}

	if !mimeMatches {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File extension (%s) does not match file content (MIME: %s)", ext, detectedMIME),
		})
		return
	}

	// 生成唯一文件名
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// 生成随机字符串
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)

	filename := fmt.Sprintf("%s_%s%s", now.Format("20060102_150405"), randomStr, ext)

	// 创建目录
	// 正文图片统一放入 uploads/images/content/<year>/<month>
	uploadDir := filepath.Join("uploads", "images", "content", year, month)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// 保存文件
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回图片 URL（使用正斜杠，适配 Web）
	imageURL := fmt.Sprintf("/uploads/images/content/%s/%s/%s", year, month, filename)

	c.JSON(http.StatusOK, gin.H{
		"url":      imageURL,
		"filename": filename,
		"size":     file.Size,
	})
}
