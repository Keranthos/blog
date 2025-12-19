package controllers

import (
	"backend/config"
	"backend/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ToggleLike 切换点赞状态（点赞/取消点赞）
func ToggleLike(c *gin.Context) {
	// 添加panic恢复
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ToggleLike panic: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}()

	// 获取用户信息（从中间件中获取）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userIDUint, ok := userIDValue.(uint)
	if !ok {
		log.Printf("无效的用户ID类型: %T, 值: %v", userIDValue, userIDValue)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	articleType := c.Param("articleType")
	articleID := c.Param("articleID")

	var articleIDUint uint
	if _, err := fmt.Sscanf(articleID, "%d", &articleIDUint); err != nil {
		log.Printf("无效的文章ID: %s, 错误: %v", articleID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	log.Printf("ToggleLike: userID=%d, articleType=%s, articleID=%d", userIDUint, articleType, articleIDUint)

	// 检查是否已经点赞（包括软删除的记录）
	var existingLike models.Like
	err := config.DB.Unscoped().Where("user_id = ? AND article_type = ? AND article_id = ?", 
		userIDUint, articleType, articleIDUint).First(&existingLike).Error

	var isLiked bool
	var likeCount int

	if err == nil {
		// 找到了记录（可能是软删除的）
		// 检查记录是否被软删除（DeletedAt.Time 不为零值表示被删除）
		if !existingLike.DeletedAt.Time.IsZero() {
			// 记录被软删除，恢复它
			log.Printf("恢复点赞: userID=%d, articleType=%s, articleID=%d", userIDUint, articleType, articleIDUint)
			if err := config.DB.Unscoped().Model(&existingLike).Update("deleted_at", nil).Error; err != nil {
				log.Printf("恢复点赞失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复点赞失败", "details": err.Error()})
				return
			}
			isLiked = true
			// 增加点赞数
			likeCount, err = increaseLikeCount(articleType, articleIDUint)
			if err != nil {
				log.Printf("增加点赞数失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败", "details": err.Error()})
				return
			}
		} else {
			// 记录存在且未删除，执行取消点赞
			log.Printf("取消点赞: userID=%d, articleType=%s, articleID=%d", userIDUint, articleType, articleIDUint)
			if err := config.DB.Delete(&existingLike).Error; err != nil {
				log.Printf("取消点赞失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败", "details": err.Error()})
				return
			}
			isLiked = false
			// 减少点赞数
			likeCount, err = decreaseLikeCount(articleType, articleIDUint)
			if err != nil {
				log.Printf("减少点赞数失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败", "details": err.Error()})
				return
			}
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 未点赞，执行点赞
		log.Printf("执行点赞: userID=%d, articleType=%s, articleID=%d", userIDUint, articleType, articleIDUint)
		// 先验证文章是否存在
		if err := validateArticleExists(articleType, articleIDUint); err != nil {
			log.Printf("文章不存在: articleType=%s, articleID=%d, 错误: %v", articleType, articleIDUint, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在", "details": err.Error()})
			return
		}
		
		newLike := models.Like{
			UserID:     userIDUint,
			ArticleType: articleType,
			ArticleID:  articleIDUint,
		}
		if err := config.DB.Create(&newLike).Error; err != nil {
			log.Printf("点赞失败: %v", err)
			// 如果是唯一索引冲突，可能是并发问题，尝试查询一次
			if strings.Contains(err.Error(), "Duplicate entry") {
				log.Printf("检测到唯一索引冲突，可能是并发问题，重新查询...")
				var conflictLike models.Like
				if queryErr := config.DB.Unscoped().Where("user_id = ? AND article_type = ? AND article_id = ?", 
					userIDUint, articleType, articleIDUint).First(&conflictLike).Error; queryErr == nil {
					// 找到了冲突的记录，恢复它（如果被软删除）
					if !conflictLike.DeletedAt.Time.IsZero() {
						if updateErr := config.DB.Unscoped().Model(&conflictLike).Update("deleted_at", nil).Error; updateErr != nil {
							log.Printf("恢复冲突记录失败: %v", updateErr)
							c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败", "details": err.Error()})
							return
						}
					}
					isLiked = true
					var countErr error
					likeCount, countErr = increaseLikeCount(articleType, articleIDUint)
					if countErr != nil {
						log.Printf("增加点赞数失败: %v", countErr)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败", "details": countErr.Error()})
						return
					}
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败", "details": err.Error()})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败", "details": err.Error()})
				return
			}
		} else {
			isLiked = true
			// 增加点赞数
			likeCount, err = increaseLikeCount(articleType, articleIDUint)
			if err != nil {
				log.Printf("增加点赞数失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败", "details": err.Error()})
				return
			}
		}
	} else {
		// 其他数据库错误
		log.Printf("查询点赞状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询点赞状态失败", "details": err.Error()})
		return
	}

	log.Printf("ToggleLike成功: isLiked=%v, likeCount=%d", isLiked, likeCount)
	c.JSON(http.StatusOK, gin.H{
		"isLiked":   isLiked,
		"likeCount": likeCount,
	})
}

// GetLikeStatus 获取文章的点赞状态和数量
func GetLikeStatus(c *gin.Context) {
	articleType := c.Param("articleType")
	articleID := c.Param("articleID")

	var articleIDUint uint
	if _, err := fmt.Sscanf(articleID, "%d", &articleIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取点赞数
	var likeCount int
	switch articleType {
	case "blog":
		var article models.BlogArticle
		if err := config.DB.First(&article, articleIDUint).Error; err == nil {
			likeCount = article.LikeCount
		}
	case "research":
		var article models.ResearchArticle
		if err := config.DB.First(&article, articleIDUint).Error; err == nil {
			likeCount = article.LikeCount
		}
	case "project":
		var article models.ProjectArticle
		if err := config.DB.First(&article, articleIDUint).Error; err == nil {
			likeCount = article.LikeCount
		}
	case "moment":
		var article models.Moment
		if err := config.DB.First(&article, articleIDUint).Error; err == nil {
			likeCount = article.LikeCount
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章类型"})
		return
	}

	// 检查用户是否已点赞（如果已登录）
	isLiked := false
	userIDValue, exists := c.Get("user_id")
	if exists {
		if userIDUint, ok := userIDValue.(uint); ok {
			var like models.Like
			if err := config.DB.Where("user_id = ? AND article_type = ? AND article_id = ?", 
				userIDUint, articleType, articleIDUint).First(&like).Error; err == nil {
				isLiked = true
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isLiked":   isLiked,
		"likeCount": likeCount,
	})
}

// increaseLikeCount 增加文章点赞数
func increaseLikeCount(articleType string, articleID uint) (int, error) {
	var likeCount int
	var err error
	
	switch articleType {
	case "blog":
		var article models.BlogArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		article.LikeCount++
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "research":
		var article models.ResearchArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		article.LikeCount++
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "project":
		var article models.ProjectArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		article.LikeCount++
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "moment":
		var article models.Moment
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		article.LikeCount++
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	default:
		return 0, fmt.Errorf("无效的文章类型: %s", articleType)
	}
	return likeCount, nil
}

// decreaseLikeCount 减少文章点赞数
func decreaseLikeCount(articleType string, articleID uint) (int, error) {
	var likeCount int
	var err error
	
	switch articleType {
	case "blog":
		var article models.BlogArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		if article.LikeCount > 0 {
			article.LikeCount--
		}
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "research":
		var article models.ResearchArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		if article.LikeCount > 0 {
			article.LikeCount--
		}
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "project":
		var article models.ProjectArticle
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		if article.LikeCount > 0 {
			article.LikeCount--
		}
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	case "moment":
		var article models.Moment
		if err = config.DB.First(&article, articleID).Error; err != nil {
			return 0, fmt.Errorf("文章不存在: %w", err)
		}
		if article.LikeCount > 0 {
			article.LikeCount--
		}
		if err = config.DB.Save(&article).Error; err != nil {
			return 0, fmt.Errorf("保存失败: %w", err)
		}
		likeCount = article.LikeCount
	default:
		return 0, fmt.Errorf("无效的文章类型: %s", articleType)
	}
	return likeCount, nil
}

// validateArticleExists 验证文章是否存在
func validateArticleExists(articleType string, articleID uint) error {
	switch articleType {
	case "blog":
		var article models.BlogArticle
		if err := config.DB.First(&article, articleID).Error; err != nil {
			return err
		}
	case "research":
		var article models.ResearchArticle
		if err := config.DB.First(&article, articleID).Error; err != nil {
			return err
		}
	case "project":
		var article models.ProjectArticle
		if err := config.DB.First(&article, articleID).Error; err != nil {
			return err
		}
	case "moment":
		var article models.Moment
		if err := config.DB.First(&article, articleID).Error; err != nil {
			return err
		}
	default:
		return fmt.Errorf("无效的文章类型: %s", articleType)
	}
	return nil
}

