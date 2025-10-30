package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取碎碎念列表
func GetMoments(c *gin.Context) {
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

	var moments []models.Moment
	offset := (page - 1) * limit

	if err := config.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&moments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get moments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": moments})
}

// 获取单个碎碎念
func GetMoment(c *gin.Context) {
	id := c.Param("id")
	var moment models.Moment

	if err := config.DB.First(&moment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Moment not found"})
		return
	}

	// 增加阅读量
	config.DB.Model(&moment).Update("view_count", moment.ViewCount+1)
	moment.ViewCount++

	c.JSON(http.StatusOK, gin.H{"data": moment})
}

// 获取碎碎念数量
func GetMomentsCount(c *gin.Context) {
	var count int64
	if err := config.DB.Model(&models.Moment{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get moments count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"num": count})
}

// 创建碎碎念
func CreateMoment(c *gin.Context) {
	var moment models.Moment
	if err := c.ShouldBindJSON(&moment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newURL, err := utils.FetchAndStoreImage(moment.Image); err == nil && newURL != "" {
		moment.Image = newURL
	}

	if err := config.DB.Create(&moment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create moment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Moment created", "data": moment})
}

// 更新碎碎念
func UpdateMoment(c *gin.Context) {
	id := c.Param("id")
	var moment models.Moment

	if err := config.DB.First(&moment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Moment not found"})
		return
	}

	var updateData models.Moment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newURL, err := utils.FetchAndStoreImage(updateData.Image); err == nil && newURL != "" {
		updateData.Image = newURL
	}

	if err := config.DB.Model(&moment).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update moment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Moment updated", "data": moment})
}

// 删除碎碎念
func DeleteMoment(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Moment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete moment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Moment deleted"})
}
