package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取文章列表
func GetArticles(c *gin.Context) {
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
	articleType := c.Query("type")
	var articles interface{}
	offset := (page - 1) * limit

	switch articleType {
	case "blog":
		var blogArticles []models.BlogArticle
		if err := config.DB.Select("id", "title", "content", "tags", "image", "view_count", "created_at").Order("created_at DESC").Offset(offset).Limit(limit).Find(&blogArticles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
		articles = blogArticles
	case "research":
		var researchArticles []models.ResearchArticle
		if err := config.DB.Select("id", "title", "abstract", "tags", "image", "view_count", "created_at").Order("created_at DESC").Offset(offset).Limit(limit).Find(&researchArticles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
		articles = researchArticles
	case "project":
		var projectArticles []models.ProjectArticle
		if err := config.DB.Select("id", "title", "status", "tags", "image", "view_count", "created_at").Order("created_at DESC").Offset(offset).Limit(limit).Find(&projectArticles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
		articles = projectArticles
	case "moment":
		var momentArticles []models.Moment
		if err := config.DB.Select("id", "title", "content", "tags", "image", "view_count", "created_at").Order("created_at DESC").Offset(offset).Limit(limit).Find(&momentArticles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
			return
		}
		articles = momentArticles
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": articles})
}

// 获取文章数量
func GetArticleCount(c *gin.Context) {
	articleType := c.Query("type")
	var count int64

	switch articleType {
	case "blog":
		if err := config.DB.Model(&models.BlogArticle{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article count"})
			return
		}
	case "research":
		if err := config.DB.Model(&models.ResearchArticle{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article count"})
			return
		}
	case "project":
		if err := config.DB.Model(&models.ProjectArticle{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article count"})
			return
		}
	case "moment":
		if err := config.DB.Model(&models.Moment{}).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article count"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"num": count})
}

// 获取文章详情
func GetArticleById(c *gin.Context) {
	id := c.Param("id")
	articleType := c.Query("type")
	var article interface{}

	switch articleType {
	case "blog":
		var blogArticle models.BlogArticle
		if err := config.DB.First(&blogArticle, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		// 增加阅读量
		config.DB.Model(&blogArticle).Update("view_count", blogArticle.ViewCount+1)
		blogArticle.ViewCount++
		article = blogArticle
	case "research":
		var researchArticle models.ResearchArticle
		if err := config.DB.First(&researchArticle, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		// 增加阅读量
		config.DB.Model(&researchArticle).Update("view_count", researchArticle.ViewCount+1)
		researchArticle.ViewCount++
		article = researchArticle
	case "project":
		var projectArticle models.ProjectArticle
		if err := config.DB.First(&projectArticle, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		// 增加阅读量
		config.DB.Model(&projectArticle).Update("view_count", projectArticle.ViewCount+1)
		projectArticle.ViewCount++
		article = projectArticle
	case "moment":
		var momentArticle models.Moment
		if err := config.DB.First(&momentArticle, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		// 增加阅读量
		config.DB.Model(&momentArticle).Update("view_count", momentArticle.ViewCount+1)
		momentArticle.ViewCount++
		article = momentArticle
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": article})
}

// 创建文章
func CreateArticle(c *gin.Context) {
	var article interface{}
	articleType := c.Query("type")

	switch articleType {
	case "blog":
		var blogArticle models.BlogArticle
		if err := c.ShouldBindJSON(&blogArticle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证文章数据
		if err := utils.ValidateArticleTitle(blogArticle.Title); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateArticleContent(blogArticle.Content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateTags(blogArticle.Tags); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateImageURL(blogArticle.Image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		article = blogArticle
	case "research":
		var researchArticle models.ResearchArticle
		if err := c.ShouldBindJSON(&researchArticle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证研究文章数据
		if err := utils.ValidateArticleTitle(researchArticle.Title); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateArticleContent(researchArticle.Abstract); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateTags(researchArticle.Tags); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateImageURL(researchArticle.Image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		article = researchArticle
	case "project":
		var projectArticle models.ProjectArticle
		if err := c.ShouldBindJSON(&projectArticle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证项目文章数据
		if err := utils.ValidateArticleTitle(projectArticle.Title); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateArticleContent(projectArticle.Status); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateTags(projectArticle.Tags); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := utils.ValidateImageURL(projectArticle.Image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		article = projectArticle
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	if err := config.DB.Create(article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article created"})
}

// 更新文章
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	articleType := c.Query("type")

	// 添加调试日志
	fmt.Printf("更新文章请求 - ID: %s, Type: %s\n", id, articleType)

	switch articleType {
	case "blog":
		var blogArticle models.BlogArticle
		if err := c.ShouldBindJSON(&blogArticle); err != nil {
			fmt.Printf("JSON绑定错误: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("博客文章数据: %+v\n", blogArticle)
		if err := config.DB.Model(&models.BlogArticle{}).Where("id = ?", id).Updates(blogArticle).Error; err != nil {
			fmt.Printf("数据库更新错误: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}
	case "research":
		var researchArticle models.ResearchArticle
		if err := c.ShouldBindJSON(&researchArticle); err != nil {
			fmt.Printf("研究文章JSON绑定错误: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("研究文章数据: %+v\n", researchArticle)
		if err := config.DB.Model(&models.ResearchArticle{}).Where("id = ?", id).Updates(researchArticle).Error; err != nil {
			fmt.Printf("研究文章数据库更新错误: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}
	case "project":
		var projectArticle models.ProjectArticle
		if err := c.ShouldBindJSON(&projectArticle); err != nil {
			fmt.Printf("项目文章JSON绑定错误: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("项目文章数据: %+v\n", projectArticle)
		if err := config.DB.Model(&models.ProjectArticle{}).Where("id = ?", id).Updates(projectArticle).Error; err != nil {
			fmt.Printf("项目文章数据库更新错误: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated"})
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	articleType := c.Query("type")

	switch articleType {
	case "blog":
		if err := config.DB.Delete(&models.BlogArticle{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
			return
		}
	case "research":
		if err := config.DB.Delete(&models.ResearchArticle{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
			return
		}
	case "project":
		if err := config.DB.Delete(&models.ProjectArticle{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}

// 获取置顶文章
func GetTopArticles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "3")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// 获取各类别的置顶文章
	var topArticles []map[string]interface{}

	// 获取置顶博客
	var blogArticles []models.BlogArticle
	if err := config.DB.Select("id", "title", "content", "tags", "image", "view_count", "created_at", "is_top").Where("is_top = ?", true).Order("created_at DESC").Find(&blogArticles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top blog articles"})
		return
	}
	for _, article := range blogArticles {
		topArticles = append(topArticles, map[string]interface{}{
			"ID":        article.ID,
			"title":     article.Title,
			"content":   article.Content,
			"tags":      article.Tags,
			"image":     article.Image,
			"viewCount": article.ViewCount,
			"CreatedAt": article.CreatedAt,
			"type":      "blog",
		})
	}

	// 获取置顶科研文章
	var researchArticles []models.ResearchArticle
	if err := config.DB.Select("id", "title", "abstract", "tags", "image", "view_count", "created_at", "is_top").Where("is_top = ?", true).Order("created_at DESC").Find(&researchArticles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top research articles"})
		return
	}
	for _, article := range researchArticles {
		topArticles = append(topArticles, map[string]interface{}{
			"ID":        article.ID,
			"title":     article.Title,
			"content":   article.Abstract,
			"tags":      article.Tags,
			"image":     article.Image,
			"viewCount": article.ViewCount,
			"CreatedAt": article.CreatedAt,
			"type":      "research",
		})
	}

	// 获取置顶项目文章
	var projectArticles []models.ProjectArticle
	if err := config.DB.Select("id", "title", "status", "tags", "image", "view_count", "created_at", "is_top").Where("is_top = ?", true).Order("created_at DESC").Find(&projectArticles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top project articles"})
		return
	}
	for _, article := range projectArticles {
		topArticles = append(topArticles, map[string]interface{}{
			"ID":        article.ID,
			"title":     article.Title,
			"content":   article.Status,
			"tags":      article.Tags,
			"image":     article.Image,
			"viewCount": article.ViewCount,
			"CreatedAt": article.CreatedAt,
			"type":      "project",
		})
	}

	// 按创建时间排序，取最新的limit篇
	// 首先按时间排序
	sort.Slice(topArticles, func(i, j int) bool {
		timeI := topArticles[i]["CreatedAt"].(time.Time)
		timeJ := topArticles[j]["CreatedAt"].(time.Time)
		return timeI.After(timeJ)
	})

	// 取最新的limit篇
	if len(topArticles) > limit {
		topArticles = topArticles[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    topArticles,
		"message": "Top articles retrieved successfully",
		"userToast": false, // 不需要看板娘提示
	})
}
