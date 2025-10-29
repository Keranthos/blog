package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"log"
	"net/http"

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

	// 获取所有评论（不包含回复），按创建时间倒序
	if err := config.DB.Where("parent_id IS NULL").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{"data": result})
}
