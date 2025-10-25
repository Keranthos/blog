package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetQuestions(c *gin.Context) {
	var questions []models.Question
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

	if err := config.DB.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证问题数据
	if err := utils.ValidateQuestionContent(question.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证作者字段
	if question.Author == "" {
		question.Author = "匿名用户"
	}

	// 检查是否包含恶意内容
	if utils.ContainsMaliciousContent(question.Content) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题内容包含不当信息"})
		return
	}

	// 清理输入内容
	question.Content = utils.SanitizeInput(question.Content)

	if err := config.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question created"})
}

func AnswerQuestion(c *gin.Context) {
	var question models.Question
	questionID := c.Param("questionId")
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证回答数据
	if err := utils.ValidateAnswerContent(question.Answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否包含恶意内容
	if utils.ContainsMaliciousContent(question.Answer) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "回答内容包含不当信息"})
		return
	}

	// 清理输入内容
	question.Answer = utils.SanitizeInput(question.Answer)

	if err := config.DB.Model(&models.Question{}).Where("id = ?", questionID).Update("answer", question.Answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to answer question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question answered"})
}
