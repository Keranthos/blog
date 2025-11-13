package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 全文搜索功能
func SearchContent(c *gin.Context) {
	keyword := c.Query("keyword")
	contentType := c.Query("type") // blog, research, project, moment, all
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	offset := (page - 1) * limit
	var results []gin.H

	// 搜索博客文章
	if contentType == "all" || contentType == "blog" {
		var blogs []models.BlogArticle
		config.DB.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Order("created_at DESC").Offset(offset).Limit(limit).Find(&blogs)

		for _, blog := range blogs {
			results = append(results, gin.H{
				"id":         blog.ID,
				"type":       "blog",
				"title":      blog.Title,
				"content":    truncateContent(blog.Content, 200),
				"created_at": blog.CreatedAt,
				"tags":       blog.Tags,
				"image":      blog.Image,
				"view_count": blog.ViewCount,
			})
		}
	}

	// 搜索科研文章
	if contentType == "all" || contentType == "research" {
		var research []models.ResearchArticle
		config.DB.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Order("created_at DESC").Offset(offset).Limit(limit).Find(&research)

		for _, article := range research {
			results = append(results, gin.H{
				"id":         article.ID,
				"type":       "research",
				"title":      article.Title,
				"content":    truncateContent(article.Content, 200),
				"created_at": article.CreatedAt,
				"tags":       article.Tags,
				"image":      article.Image,
				"view_count": article.ViewCount,
			})
		}
	}

	// 搜索项目文章
	if contentType == "all" || contentType == "project" {
		var projects []models.ProjectArticle
		config.DB.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Order("created_at DESC").Offset(offset).Limit(limit).Find(&projects)

		for _, project := range projects {
			results = append(results, gin.H{
				"id":         project.ID,
				"type":       "project",
				"title":      project.Title,
				"content":    truncateContent(project.Content, 200),
				"created_at": project.CreatedAt,
				"tags":       project.Tags,
				"image":      project.Image,
				"view_count": project.ViewCount,
			})
		}
	}

	// 搜索碎碎念
	if contentType == "all" || contentType == "moment" {
		var moments []models.Moment
		config.DB.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Order("created_at DESC").Offset(offset).Limit(limit).Find(&moments)

		for _, moment := range moments {
			results = append(results, gin.H{
				"id":         moment.ID,
				"type":       "moment",
				"title":      moment.Title,
				"content":    truncateContent(moment.Content, 200),
				"created_at": moment.CreatedAt,
				"tags":       []string{},
				"image":      moment.Image,
				"view_count": moment.ViewCount,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": len(results),
		"page":  page,
		"limit": limit,
	})
}

// 截断内容为指定长度
func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}
