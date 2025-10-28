package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type PresentationController struct {
	DB *gorm.DB
}

func NewPresentationController(db *gorm.DB) *PresentationController {
	return &PresentationController{DB: db}
}

// CreatePresentation 创建讲演
func (pc *PresentationController) CreatePresentation(c *gin.Context) {
	// 获取用户ID（从JWT token中解析）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 解析表单数据
	title := c.PostForm("title")
	tags := c.PostForm("tags")
	isPrivateStr := c.PostForm("is_private")

	// 解析是否私密
	isPrivate := false
	if isPrivateStr == "true" || isPrivateStr == "1" {
		isPrivate = true
	}

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题不能为空"})
		return
	}

	// 处理文件上传
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败: " + err.Error()})
		return
	}
	defer file.Close()

	// 验证文件类型
	allowedTypes := []string{".pdf"}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !contains(allowedTypes, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅支持 .pdf"})
		return
	}

	// 验证文件大小 (限制为50MB)
	if header.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过50MB"})
		return
	}

	// 创建上传目录
	uploadDir := "uploads/presentations"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败"})
		return
	}

	// 处理缩略图
	var thumbnailPath string

	// 优先处理上传的缩略图文件
	if thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail"); err == nil {
		defer thumbnailFile.Close()

		// 验证缩略图类型
		if strings.HasPrefix(thumbnailHeader.Header.Get("Content-Type"), "image/") {
			thumbnailFilename := fmt.Sprintf("%d_thumbnail_%s", timestamp, thumbnailHeader.Filename)
			thumbnailPath = filepath.Join(uploadDir, thumbnailFilename)

			thumbnailDst, err := os.Create(thumbnailPath)
			if err == nil {
				io.Copy(thumbnailDst, thumbnailFile)
				thumbnailDst.Close()
			}
		}
	} else {
		// 如果没有上传文件，检查是否有缩略图URL
		if thumbnailURL := c.PostForm("thumbnail_url"); thumbnailURL != "" {
			thumbnailPath = thumbnailURL
		}
	}

	// 创建讲演记录
	presentation := models.Presentation{
		Title:     title,
		Tags:      tags,
		FilePath:  filePath,
		Thumbnail: thumbnailPath,
		FileSize:  header.Size,
		FileType:  strings.TrimPrefix(ext, "."),
		UserID:    userID.(uint),
		IsActive:  true,
		IsPrivate: isPrivate,
		Password:  "0714", // 固定密码
	}

	if err := pc.DB.Create(&presentation).Error; err != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(filePath)
		if thumbnailPath != "" {
			os.Remove(thumbnailPath)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存讲演信息失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "讲演创建成功",
		"data":    presentation.ToResponse(),
	})
}

// GetPresentations 获取讲演列表
func (pc *PresentationController) GetPresentations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "9"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 9
	}

	offset := (page - 1) * pageSize

	var presentations []models.Presentation
	var total int64

	// 查询总数
	pc.DB.Model(&models.Presentation{}).Where("is_active = ?", true).Count(&total)

	// 查询数据
	if err := pc.DB.Where("is_active = ?", true).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&presentations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演列表失败"})
		return
	}

	// 转换为响应格式
	var responses []models.PresentationResponse
	for _, p := range presentations {
		responses = append(responses, p.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetPresentation 获取单个讲演详情
func (pc *PresentationController) GetPresentation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Where("id = ? AND is_active = ?", id, true).
		Preload("User").
		First(&presentation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演失败"})
		}
		return
	}

	// 增加浏览次数
	pc.DB.Model(&presentation).Update("view_count", gorm.Expr("view_count + 1"))

	c.JSON(http.StatusOK, gin.H{
		"data": presentation.ToResponse(),
	})
}

// UpdatePresentation 更新讲演
func (pc *PresentationController) UpdatePresentation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&presentation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在或无权限修改"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演失败"})
		}
		return
	}

	// 解析更新数据
	var updateData models.PresentationUpdateRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if updateData.Title != "" {
		updates["title"] = updateData.Title
	}
	if updateData.Tags != "" {
		updates["tags"] = updateData.Tags
	}
	if updateData.IsActive != nil {
		updates["is_active"] = *updateData.IsActive
	}
	if updateData.IsPrivate != nil {
		updates["is_private"] = *updateData.IsPrivate
	}

	if err := pc.DB.Model(&presentation).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新讲演失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "讲演更新成功",
		"data":    presentation.ToResponse(),
	})
}

// DeletePresentation 删除讲演
func (pc *PresentationController) DeletePresentation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&presentation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在或无权限删除"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演失败"})
		}
		return
	}

	// 软删除（设置为非激活状态）
	if err := pc.DB.Model(&presentation).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除讲演失败"})
		return
	}

	// 可选：删除文件
	// os.Remove(presentation.FilePath)
	// if presentation.Thumbnail != "" {
	//     os.Remove(presentation.Thumbnail)
	// }

	c.JSON(http.StatusOK, gin.H{"message": "讲演删除成功"})
}

// ServePresentationFile 提供讲演文件下载
func (pc *PresentationController) ServePresentationFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Where("id = ? AND is_active = ?", id, true).First(&presentation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演失败"})
		}
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(presentation.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(presentation.FilePath)))

	// 根据文件类型设置正确的Content-Type
	switch presentation.FileType {
	case "pdf":
		c.Header("Content-Type", "application/pdf")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	// 提供文件下载
	c.File(presentation.FilePath)
}

// PreviewPresentationFile 提供讲演文件预览（用于iframe）
func (pc *PresentationController) PreviewPresentationFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Preload("User").First(&presentation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(presentation.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置预览响应头（不设置Content-Disposition，让浏览器直接显示）
	switch presentation.FileType {
	case "pdf":
		c.Header("Content-Type", "application/pdf")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	// 提供文件预览
	c.File(presentation.FilePath)
}

// GetPresentationInfo 获取讲演文件信息（包括PDF页数等）
func (pc *PresentationController) GetPresentationInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Preload("User").First(&presentation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(presentation.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 对于PDF文件，尝试获取页数信息
	var totalPages int = 1
	if presentation.FileType == "pdf" {
		// 这里可以添加PDF页数解析逻辑
		// 目前返回默认值，实际项目中可以使用PDF库来解析
		totalPages = 10 // 默认值
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          presentation.ID,
		"title":       presentation.Title,
		"file_type":   presentation.FileType,
		"file_size":   presentation.FileSize,
		"total_pages": totalPages,
		"is_private":  presentation.IsPrivate,
	})
}

// VerifyPresentationPassword 验证讲演密码
func (pc *PresentationController) VerifyPresentationPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var presentation models.Presentation
	if err := pc.DB.Where("id = ? AND is_active = ?", id, true).First(&presentation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "讲演不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询讲演失败"})
		}
		return
	}

	// 如果不是私密讲演，直接返回成功
	if !presentation.IsPrivate {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "讲演为公开状态"})
		return
	}

	// 解析密码
	var passwordData struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 验证密码
	if passwordData.Password == presentation.Password {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "密码验证成功"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
	}
}

// 辅助函数
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
