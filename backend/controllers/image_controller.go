package controllers

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ImageInfo 图片信息
type ImageInfo struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadTime int64  `json:"uploadTime"`
}

// GetImages 获取所有上传的图片列表
func GetImages(c *gin.Context) {
	images := []ImageInfo{}

	// 遍历 uploads/images 目录
	uploadsDir := "uploads/images"
	err := filepath.WalkDir(uploadsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理图片文件
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
			info, err := d.Info()
			if err != nil {
				return err
			}

			// 转换为 Web 路径（使用正斜杠）
			webPath := "/" + strings.ReplaceAll(path, "\\", "/")

			images = append(images, ImageInfo{
				Path:       webPath,
				Name:       d.Name(),
				Size:       info.Size(),
				UploadTime: info.ModTime().Unix() * 1000, // 转换为毫秒
			})
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read images directory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// DeleteImage 删除图片
func DeleteImage(c *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 安全检查：确保路径在 uploads 目录内
	if !strings.HasPrefix(req.Path, "/uploads/images/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	// 移除开头的斜杠，转换为文件系统路径
	filePath := strings.TrimPrefix(req.Path, "/")
	filePath = filepath.FromSlash(filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// 尝试删除空的父目录
	dir := filepath.Dir(filePath)
	for strings.HasPrefix(dir, "uploads"+string(filepath.Separator)+"images") {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}
		os.Remove(dir)
		dir = filepath.Dir(dir)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// RandomImageRequest 随机图片请求
type RandomImageRequest struct {
	APIUrl  string `json:"apiUrl"`
	APIName string `json:"apiName"`
}

// RandomImageResponse 随机图片响应
type RandomImageResponse struct {
	Success  bool   `json:"success"`
	ImageURL string `json:"imageUrl,omitempty"`
	Message  string `json:"message,omitempty"`
}

// GetRandomImage 获取随机图片
func GetRandomImage(c *gin.Context) {
	var req RandomImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("绑定请求失败: %v\n", err)
		c.JSON(http.StatusBadRequest, RandomImageResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	fmt.Printf("收到随机图片请求\n")

	// 只使用DMoe API
	dmoeURL := "https://www.dmoe.cc/random.php"
	fmt.Printf("使用DMoe API: %s\n", dmoeURL)

	c.JSON(http.StatusOK, RandomImageResponse{
		Success:  true,
		ImageURL: dmoeURL,
	})
}
