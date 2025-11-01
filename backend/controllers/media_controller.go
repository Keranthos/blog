package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取媒体列表
func GetMedia(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	// 分页上限校验：防止 limit 过大导致性能问题
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}
	mediaType := c.Query("type")
	if mediaType != "novels" && mediaType != "books" && mediaType != "movies" {
		c.JSON(http.StatusBadRequest, gin.H{"error": mediaType + " is not a valid media type"})
		return
	}

	var media []models.Media
	offset := (page - 1) * limit

	if err := config.DB.Where("type = ?", mediaType).Order("created_at DESC").Offset(offset).Limit(limit).Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"medias": media,
		},
	})
}

// 获取单个媒体
func GetMediaByID(c *gin.Context) {
	mediaID := c.Param("mediaId")
	mediaType := c.Query("type")

	if mediaType != "novels" && mediaType != "books" && mediaType != "movies" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type"})
		return
	}

	var media models.Media
	if err := config.DB.Where("id = ? AND type = ?", mediaID, mediaType).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": media,
	})
}

func CreateMedia(c *gin.Context) {
	var media models.Media
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 本地化外链封面到 uploads/images/media
	if media.Poster != "" {
		if local, err := utils.FetchAndStoreImageTo(media.Poster, "media"); err == nil {
			media.Poster = local
		} else {
			// 下载失败不阻断创建，仅记录到响应
			c.Set("media_poster_error", err.Error())
		}
	}

	if err := config.DB.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media created"})
}

func UpdateMedia(c *gin.Context) {
	var media models.Media
	mediaID := c.Param("mediaId")
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 本地化外链封面到 uploads/images/media
	if media.Poster != "" {
		if local, err := utils.FetchAndStoreImageTo(media.Poster, "media"); err == nil {
			media.Poster = local
		} else {
			c.Set("media_poster_error", err.Error())
		}
	}

	if err := config.DB.Model(&models.Media{}).Where("id = ?", mediaID).Updates(media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media updated"})
}

func DeleteMedia(c *gin.Context) {
	mediaID := c.Param("mediaId")

	// 先获取媒体信息，以便删除封面图片
	var media models.Media
	if err := config.DB.Where("id = ?", mediaID).First(&media).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// 删除封面图片
	if media.Poster != "" {
		utils.DeleteImageFileSafely(media.Poster)
	}

	// 从评论内容中提取并删除所有图片（如果有Review字段包含图片）
	if media.Review != "" {
		imageURLs := utils.ExtractImageURLsFromContent(media.Review)
		for _, imgURL := range imageURLs {
			utils.DeleteImageFileSafely(imgURL)
		}
	}

	// 删除媒体记录
	if err := config.DB.Where("id = ?", mediaID).Delete(&models.Media{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted"})
}

// 批量本地化历史媒体封面
func MigrateMediaImages(c *gin.Context) {
	var items []models.Media
	if err := config.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}

	migrated := 0
	failed := 0
	skipped := 0

	for i := range items {
		m := &items[i]
		if m.Poster == "" {
			skipped++
			continue
		}

		// 只处理外链图片，本地图片跳过
		if !utils.IsExternalImageURL(m.Poster) {
			skipped++
			continue
		}

		// 本地化外链图片
		local, err := utils.FetchAndStoreImageTo(m.Poster, "media")
		if err != nil {
			failed++
			continue
		}

		// 只有当返回的路径与原始路径不同时才更新（说明已经下载并保存）
		if local != m.Poster {
			if err := config.DB.Model(&models.Media{}).Where("id = ?", m.ID).Update("poster", local).Error; err == nil {
				migrated++
			} else {
				failed++
			}
		} else {
			skipped++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"migrated": migrated,
		"failed":   failed,
		"skipped":  skipped,
		"total":    len(items),
		"message":  fmt.Sprintf("迁移完成：成功 %d 个，失败 %d 个，跳过 %d 个", migrated, failed, skipped),
	})
}
