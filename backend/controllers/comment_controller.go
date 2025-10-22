package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetComments(c *gin.Context) {
	var comments []models.Comment
	blogID := c.Param("blogID")
	commentType := c.Query("type")

	// 只获取顶级评论（ParentID 为 null）
	if err := config.DB.Where("blog_id = ? AND type = ? AND parent_id IS NULL", blogID, commentType).Preload("Replies").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证评论数据
	if err := utils.ValidateCommentContent(comment.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否包含恶意内容
	if utils.ContainsMaliciousContent(comment.Content) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容包含不当信息"})
		return
	}

	// 清理输入内容
	comment.Content = utils.SanitizeInput(comment.Content)

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created"})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	if err := config.DB.Delete(&models.Comment{}, commentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
