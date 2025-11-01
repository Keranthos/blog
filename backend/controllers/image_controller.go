package controllers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/config"
	"backend/models"

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

// ImageUsage 图片使用位置信息
type ImageUsage struct {
	Type      string `json:"type"`                // "blog", "research", "project", "moment", "media"
	ID        uint   `json:"id"`                  // 文章/媒体ID
	Title     string `json:"title"`               // 标题
	UsageType string `json:"usageType"`           // "cover" 封面, "content" 正文
	MediaType string `json:"mediaType,omitempty"` // 仅当type为"media"时，值为"novels"/"books"/"movies"
}

// checkImageInUse 检查图片是否正在被使用（返回简单描述）
func checkImageInUse(imagePath string) ([]string, error) {
	var usageList []string

	// 转换为Web路径格式（统一使用正斜杠）
	webPath := "/" + strings.ReplaceAll(imagePath, "\\", "/")

	// 检查博客文章封面
	var blogCount int64
	config.DB.Model(&models.BlogArticle{}).Where("image = ?", webPath).Count(&blogCount)
	if blogCount > 0 {
		usageList = append(usageList, fmt.Sprintf("博客文章封面 (%d篇)", blogCount))
	}

	// 检查博客文章正文内容
	var blogContentCount int64
	config.DB.Model(&models.BlogArticle{}).Where("content LIKE ?", "%"+webPath+"%").Count(&blogContentCount)
	if blogContentCount > 0 {
		usageList = append(usageList, fmt.Sprintf("博客文章正文 (%d篇)", blogContentCount))
	}

	// 检查科研文章封面
	var researchCount int64
	config.DB.Model(&models.ResearchArticle{}).Where("image = ?", webPath).Count(&researchCount)
	if researchCount > 0 {
		usageList = append(usageList, fmt.Sprintf("科研文章封面 (%d篇)", researchCount))
	}

	// 检查科研文章正文（如果有内容字段）
	// 注意：ResearchArticle可能没有Content字段，这里先检查Abstract
	var researchAbstractCount int64
	config.DB.Model(&models.ResearchArticle{}).Where("abstract LIKE ?", "%"+webPath+"%").Count(&researchAbstractCount)
	if researchAbstractCount > 0 {
		usageList = append(usageList, fmt.Sprintf("科研文章摘要 (%d篇)", researchAbstractCount))
	}

	// 检查项目文章封面
	var projectCount int64
	config.DB.Model(&models.ProjectArticle{}).Where("image = ?", webPath).Count(&projectCount)
	if projectCount > 0 {
		usageList = append(usageList, fmt.Sprintf("项目文章封面 (%d篇)", projectCount))
	}

	// 检查媒体（书影集）封面
	var mediaCount int64
	config.DB.Model(&models.Media{}).Where("poster = ?", webPath).Count(&mediaCount)
	if mediaCount > 0 {
		usageList = append(usageList, fmt.Sprintf("书影集封面 (%d个)", mediaCount))
	}

	// 检查媒体（书影集）正文内容
	var mediaReviewCount int64
	config.DB.Model(&models.Media{}).Where("review LIKE ?", "%"+webPath+"%").Count(&mediaReviewCount)
	if mediaReviewCount > 0 {
		usageList = append(usageList, fmt.Sprintf("书影集正文 (%d个)", mediaReviewCount))
	}

	// 检查随笔封面
	var momentCount int64
	config.DB.Model(&models.Moment{}).Where("image = ?", webPath).Count(&momentCount)
	if momentCount > 0 {
		usageList = append(usageList, fmt.Sprintf("随笔封面 (%d篇)", momentCount))
	}

	// 检查随笔内容
	var momentContentCount int64
	config.DB.Model(&models.Moment{}).Where("content LIKE ?", "%"+webPath+"%").Count(&momentContentCount)
	if momentContentCount > 0 {
		usageList = append(usageList, fmt.Sprintf("随笔正文 (%d篇)", momentContentCount))
	}

	return usageList, nil
}

// GetImageUsages 获取图片的所有使用位置（返回详细信息）
func GetImageUsages(c *gin.Context) {
	imagePath := c.Query("path")
	if imagePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image path is required"})
		return
	}

	// 转换为Web路径格式（统一使用正斜杠）
	// 先去掉反斜杠，然后确保以 / 开头
	webPath := strings.ReplaceAll(imagePath, "\\", "/")
	if !strings.HasPrefix(webPath, "/") {
		webPath = "/" + webPath
	}

	var usages []ImageUsage

	// 调试日志：打印查询的路径
	fmt.Printf("查询图片路径: %s\n", webPath)

	// 调试：查询一个博客文章的image字段示例，查看实际存储格式
	var sampleBlog models.BlogArticle
	if err := config.DB.First(&sampleBlog).Error; err == nil {
		fmt.Printf("数据库中博客文章image字段示例: %s\n", sampleBlog.Image)
	}

	// 检查博客文章封面
	var blogArticles []models.BlogArticle
	// 尝试多种路径格式进行匹配
	pathVariations := []string{
		webPath,
		strings.TrimPrefix(webPath, "/"),
		strings.TrimSuffix(webPath, "/"),
		strings.Trim(strings.TrimPrefix(webPath, "/"), "/"),
	}

	// 构建查询条件：尝试精确匹配多种路径格式
	query := config.DB
	for i, path := range pathVariations {
		if i == 0 {
			query = query.Where("image = ?", path)
		} else {
			query = query.Or("image = ?", path)
		}
	}
	query.Find(&blogArticles)

	// 如果还是找不到，尝试LIKE模糊匹配
	if len(blogArticles) == 0 {
		// 提取文件名（路径的最后一部分）
		fileName := filepath.Base(webPath)
		config.DB.Where("image LIKE ?", "%"+fileName+"%").Find(&blogArticles)
	}

	fmt.Printf("找到博客文章: %d 篇\n", len(blogArticles))
	for _, article := range blogArticles {
		usages = append(usages, ImageUsage{
			Type:      "blog",
			ID:        article.ID,
			Title:     article.Title,
			UsageType: "cover",
		})
	}

	// 检查博客文章正文
	var blogContentArticles []models.BlogArticle
	config.DB.Where("content LIKE ?", "%"+webPath+"%").Find(&blogContentArticles)
	for _, article := range blogContentArticles {
		// 避免重复（如果封面和正文都是同一图片）
		exists := false
		for _, u := range usages {
			if u.Type == "blog" && u.ID == article.ID {
				exists = true
				break
			}
		}
		if !exists {
			usages = append(usages, ImageUsage{
				Type:      "blog",
				ID:        article.ID,
				Title:     article.Title,
				UsageType: "content",
			})
		}
	}

	// 检查科研文章封面
	var researchArticles []models.ResearchArticle
	pathVariations = []string{
		webPath,
		strings.TrimPrefix(webPath, "/"),
		strings.TrimSuffix(webPath, "/"),
		strings.Trim(strings.TrimPrefix(webPath, "/"), "/"),
	}
	query = config.DB
	for i, path := range pathVariations {
		if i == 0 {
			query = query.Where("image = ?", path)
		} else {
			query = query.Or("image = ?", path)
		}
	}
	query.Find(&researchArticles)
	if len(researchArticles) == 0 {
		fileName := filepath.Base(webPath)
		config.DB.Where("image LIKE ?", "%"+fileName+"%").Find(&researchArticles)
	}
	for _, article := range researchArticles {
		usages = append(usages, ImageUsage{
			Type:      "research",
			ID:        article.ID,
			Title:     article.Title,
			UsageType: "cover",
		})
	}

	// 检查科研文章摘要
	config.DB.Where("abstract LIKE ?", "%"+webPath+"%").Find(&researchArticles)
	for _, article := range researchArticles {
		exists := false
		for _, u := range usages {
			if u.Type == "research" && u.ID == article.ID {
				exists = true
				break
			}
		}
		if !exists {
			usages = append(usages, ImageUsage{
				Type:      "research",
				ID:        article.ID,
				Title:     article.Title,
				UsageType: "content",
			})
		}
	}

	// 检查项目文章封面
	var projectArticles []models.ProjectArticle
	pathVariations = []string{
		webPath,
		strings.TrimPrefix(webPath, "/"),
		strings.TrimSuffix(webPath, "/"),
		strings.Trim(strings.TrimPrefix(webPath, "/"), "/"),
	}
	query = config.DB
	for i, path := range pathVariations {
		if i == 0 {
			query = query.Where("image = ?", path)
		} else {
			query = query.Or("image = ?", path)
		}
	}
	query.Find(&projectArticles)
	if len(projectArticles) == 0 {
		fileName := filepath.Base(webPath)
		config.DB.Where("image LIKE ?", "%"+fileName+"%").Find(&projectArticles)
	}
	for _, article := range projectArticles {
		usages = append(usages, ImageUsage{
			Type:      "project",
			ID:        article.ID,
			Title:     article.Title,
			UsageType: "cover",
		})
	}

	// 检查媒体（书影集）封面
	var mediaItems []models.Media
	pathVariations = []string{
		webPath,
		strings.TrimPrefix(webPath, "/"),
		strings.TrimSuffix(webPath, "/"),
		strings.Trim(strings.TrimPrefix(webPath, "/"), "/"),
	}
	query = config.DB
	for i, path := range pathVariations {
		if i == 0 {
			query = query.Where("poster = ?", path)
		} else {
			query = query.Or("poster = ?", path)
		}
	}
	query.Find(&mediaItems)
	if len(mediaItems) == 0 {
		fileName := filepath.Base(webPath)
		config.DB.Where("poster LIKE ?", "%"+fileName+"%").Find(&mediaItems)
	}
	for _, media := range mediaItems {
		usages = append(usages, ImageUsage{
			Type:      "media",
			ID:        media.ID,
			Title:     media.Name,
			UsageType: "cover",
			MediaType: media.Type, // "novels"/"books"/"movies"
		})
	}

	// 检查媒体（书影集）正文内容
	config.DB.Where("review LIKE ?", "%"+webPath+"%").Find(&mediaItems)
	for _, media := range mediaItems {
		exists := false
		for _, u := range usages {
			if u.Type == "media" && u.ID == media.ID {
				exists = true
				break
			}
		}
		if !exists {
			usages = append(usages, ImageUsage{
				Type:      "media",
				ID:        media.ID,
				Title:     media.Name,
				UsageType: "content",
				MediaType: media.Type, // "novels"/"books"/"movies"
			})
		}
	}

	// 检查随笔封面
	var moments []models.Moment
	pathVariations = []string{
		webPath,
		strings.TrimPrefix(webPath, "/"),
		strings.TrimSuffix(webPath, "/"),
		strings.Trim(strings.TrimPrefix(webPath, "/"), "/"),
	}
	query = config.DB
	for i, path := range pathVariations {
		if i == 0 {
			query = query.Where("image = ?", path)
		} else {
			query = query.Or("image = ?", path)
		}
	}
	query.Find(&moments)
	if len(moments) == 0 {
		fileName := filepath.Base(webPath)
		config.DB.Where("image LIKE ?", "%"+fileName+"%").Find(&moments)
	}
	for _, moment := range moments {
		usages = append(usages, ImageUsage{
			Type:      "moment",
			ID:        moment.ID,
			Title:     moment.Title,
			UsageType: "cover",
		})
	}

	// 检查随笔内容
	config.DB.Where("content LIKE ?", "%"+webPath+"%").Find(&moments)
	for _, moment := range moments {
		exists := false
		for _, u := range usages {
			if u.Type == "moment" && u.ID == moment.ID {
				exists = true
				break
			}
		}
		if !exists {
			usages = append(usages, ImageUsage{
				Type:      "moment",
				ID:        moment.ID,
				Title:     moment.Title,
				UsageType: "content",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"usages": usages,
		"count":  len(usages),
	})
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

	// 检查图片是否正在被使用
	usageList, err := checkImageInUse(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check image usage"})
		return
	}

	if len(usageList) > 0 {
		usageMsg := "该图片正在被使用，无法删除：\n" + strings.Join(usageList, "\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Image is in use",
			"message":   usageMsg,
			"usageList": usageList,
			"userToast": true, // 需要看板娘提示
		})
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

	c.JSON(http.StatusOK, gin.H{
		"message":   "Image deleted successfully",
		"userToast": true, // 需要看板娘提示
	})
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
	APIName  string `json:"apiName,omitempty"`
	Message  string `json:"message,omitempty"`
}

// ImageAPI 图片API配置
type ImageAPI struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
}

// GetAvailableAPIs 获取可用的图片API列表
func GetAvailableAPIs(c *gin.Context) {
	apis := []ImageAPI{
		{
			Name:        "樱花动漫",
			Description: "DMoe动漫图片API",
			URL:         "https://www.dmoe.cc/random.php",
			Category:    "动漫",
		},
		{
			Name:        "ImgAPI动漫",
			Description: "ImgAPI动漫图片",
			URL:         "https://imgapi.cn/api.php?fl=dongman&gs=json",
			Category:    "动漫",
		},
		{
			Name:        "ImgAPI美女",
			Description: "ImgAPI美女图片",
			URL:         "https://imgapi.cn/api.php?fl=meizi&gs=json",
			Category:    "美女",
		},
		{
			Name:        "ImgAPI风景",
			Description: "ImgAPI风景图片",
			URL:         "https://imgapi.cn/api.php?fl=fengjing&gs=json",
			Category:    "风景",
		},
		{
			Name:        "ImgAPI随机",
			Description: "ImgAPI随机图片",
			URL:         "https://imgapi.cn/api.php?fl=suiji&gs=json",
			Category:    "随机",
		},
		{
			Name:        "CosPlay图片",
			Description: "CosPlay图片API",
			URL:         "https://imgapi.cn/cos.php?return=json",
			Category:    "CosPlay",
		},
		{
			Name:        "Bing每日一图",
			Description: "Bing每日一图API",
			URL:         "https://imgapi.cn/bing.php?rand=true",
			Category:    "风景",
		},
		{
			Name:        "懒加载图片",
			Description: "随机懒加载图片",
			URL:         "https://imgapi.cn/loading.php?return=json",
			Category:    "随机",
		},
	}

	c.JSON(http.StatusOK, gin.H{"apis": apis})
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

	// 随机选择一个API
	apis := []struct {
		name    string
		url     string
		handler func(string) (string, error)
	}{
		{"樱花动漫", "https://www.dmoe.cc/random.php", nil},
		{"ImgAPI动漫", "https://imgapi.cn/api.php?fl=dongman&gs=json", getImgAPIImage},
		{"ImgAPI美女", "https://imgapi.cn/api.php?fl=meizi&gs=json", getImgAPIImage},
		{"ImgAPI风景", "https://imgapi.cn/api.php?fl=fengjing&gs=json", getImgAPIImage},
		{"ImgAPI随机", "https://imgapi.cn/api.php?fl=suiji&gs=json", getImgAPIImage},
		{"CosPlay图片", "https://imgapi.cn/cos.php?return=json", getImgAPIImage},
		{"Bing每日一图", "https://imgapi.cn/bing.php?rand=true", func(string) (string, error) { return getBingImage() }},
		{"懒加载图片", "https://imgapi.cn/loading.php?return=json", getImgAPIImage},
	}

	// 随机选择一个API
	rand.Seed(time.Now().UnixNano())
	selectedAPI := apis[rand.Intn(len(apis))]
	fmt.Printf("随机选择API: %s\n", selectedAPI.name)

	var imageURL string
	var err error

	if selectedAPI.handler != nil {
		imageURL, err = selectedAPI.handler(selectedAPI.url)
	} else {
		// 对于樱花动漫API，直接使用URL
		imageURL = selectedAPI.url
	}

	if err != nil {
		fmt.Printf("获取图片失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, RandomImageResponse{
			Success: false,
			Message: "Failed to get image",
		})
		return
	}

	// 添加缓存破坏参数
	separator := "?"
	if strings.Contains(imageURL, "?") {
		separator = "&"
	}
	imageURL = fmt.Sprintf("%s%st=%d", imageURL, separator, time.Now().Unix())

	fmt.Printf("返回图片URL: %s\n", imageURL)

	c.JSON(http.StatusOK, RandomImageResponse{
		Success:  true,
		ImageURL: imageURL,
		APIName:  selectedAPI.name,
	})
}

// getImgAPIImage 从ImgAPI获取图片
func getImgAPIImage(apiURL string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// 根据不同的API返回格式解析图片URL
	if imgURL, ok := result["imgurl"].(string); ok {
		return imgURL, nil
	}
	if imgURL, ok := result["url"].(string); ok {
		return imgURL, nil
	}
	if imgURL, ok := result["image"].(string); ok {
		return imgURL, nil
	}

	return "", fmt.Errorf("无法从API响应中获取图片URL")
}

// getBingImage 获取Bing每日一图
func getBingImage() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://imgapi.cn/bing.php?rand=true")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if imgURL, ok := result["url"].(string); ok {
		return imgURL, nil
	}

	return "", fmt.Errorf("无法从Bing API获取图片URL")
}

// CustomImageUpload 自定义图片上传
func CustomImageUpload(c *gin.Context) {
	fmt.Printf("收到上传请求，Content-Type: %s\n", c.GetHeader("Content-Type"))
	fmt.Printf("请求方法: %s\n", c.Request.Method)

	// 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		fmt.Printf("获取文件失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无法获取上传文件",
		})
		return
	}
	defer file.Close()

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "只支持图片文件",
		})
		return
	}

	// 检查文件大小 (限制为10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件大小不能超过10MB",
		})
		return
	}

	// 创建上传目录
	uploadDir := "uploads/custom-images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建上传目录失败",
		})
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存文件失败",
		})
		return
	}

	// 返回文件URL（使用完整URL）
	webPath := "/" + strings.ReplaceAll(filePath, "\\", "/")
	fullURL := "http://localhost:3000" + webPath

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"imageUrl": fullURL,
		"fileName": header.Filename,
		"fileSize": header.Size,
		"message":  "上传成功",
	})
}
