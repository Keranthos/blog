package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetComments(c *gin.Context) {
	var comments []models.Comment
	blogID := c.Param("blogID")
	commentType := c.Query("type")

	// 只获取顶级评论（ParentID 为 null），按创建时间倒序排列（最新的在前）
	if err := config.DB.Where("blog_id = ? AND type = ? AND parent_id IS NULL", blogID, commentType).
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			// 对回复也按创建时间倒序排列
			return db.Order("created_at DESC")
		}).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Printf("CreateComment - JSON绑定错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("CreateComment - 接收到的评论: %+v", comment)

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
	
	// 如果有quoted_text，使用专门的清理函数（保留换行符）
	if comment.QuotedText != nil && *comment.QuotedText != "" {
		sanitized := utils.SanitizeQuotedText(*comment.QuotedText)
		comment.QuotedText = &sanitized
	}

	// 如果有parent_id，验证父评论是否存在
	if comment.ParentID != nil {
		var parentComment models.Comment
		if err := config.DB.First(&parentComment, *comment.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不存在"})
			return
		}
		// 确保父评论和当前评论属于同一篇文章
		if parentComment.BlogID != comment.BlogID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "回复的评论不属于当前文章"})
			return
		}
	}

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

// GetAllComments 获取所有评论（按时间倒序，包含文章信息）
func GetAllComments(c *gin.Context) {
	var comments []models.Comment
	var total int64

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// 分页上限校验
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	// 获取所有评论（包括回复），按创建时间倒序
	query := config.DB.Order("created_at DESC")
	
	// 计算总数
	if err := query.Model(&models.Comment{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
		return
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// 为每个评论查找对应的文章信息
	type CommentWithArticle struct {
		models.Comment
		ArticleTitle string `json:"articleTitle"`
		ArticleType  string `json:"articleType"`
	}

	var result []CommentWithArticle
	for _, comment := range comments {
		item := CommentWithArticle{
			Comment: comment,
		}

		// 根据评论类型和blogID查找对应的文章
		switch comment.Type {
		case "blog":
			var article models.BlogArticle
			if err := config.DB.First(&article, comment.BlogID).Error; err == nil {
				item.ArticleTitle = article.Title
				item.ArticleType = "blog"
			}
		case "moment":
			var article models.Moment
			if err := config.DB.First(&article, comment.BlogID).Error; err == nil {
				item.ArticleTitle = article.Title
				item.ArticleType = "moment"
			}
		case "research":
			var article models.ResearchArticle
			if err := config.DB.First(&article, comment.BlogID).Error; err == nil {
				item.ArticleTitle = article.Title
				item.ArticleType = "research"
			}
		case "project":
			var article models.ProjectArticle
			if err := config.DB.First(&article, comment.BlogID).Error; err == nil {
				item.ArticleTitle = article.Title
				item.ArticleType = "project"
			}
		}

		result = append(result, item)
	}

	// 计算分页信息
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasMore := page < totalPages

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"pagination": gin.H{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"hasMore":  hasMore,
		},
	})
}
