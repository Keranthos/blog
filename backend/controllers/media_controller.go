package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
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
	for i := range items {
		m := &items[i]
		if local, err := utils.FetchAndStoreImageTo(m.Poster, "media"); err == nil && local != m.Poster {
			if err := config.DB.Model(&models.Media{}).Where("id = ?", m.ID).Update("poster", local).Error; err == nil {
				migrated++
			} else {
				failed++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"migrated": migrated,
		"failed":   failed,
		"total":    len(items),
	})
}
