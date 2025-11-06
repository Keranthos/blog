package controllers

import (
	"backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlannerController 计划控制器
type PlannerController struct {
	DB *gorm.DB
}

// NewPlannerController 创建新的计划控制器
func NewPlannerController(db *gorm.DB) *PlannerController {
	return &PlannerController{DB: db}
}

// ==================== 任务相关 ====================

// GetTasks 获取任务列表
func (pc *PlannerController) GetTasks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	var tasks []models.Task
	query := pc.DB.Where("user_id = ?", userID)

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
		}
	}

	if err := query.Order("priority DESC, created_at ASC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// CreateTask 创建任务
func (pc *PlannerController) CreateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var taskData struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Priority      int    `json:"priority"`
		EstimatedTime int    `json:"estimatedTime"`
		Date          string `json:"date"`
		Completed     bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&taskData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	task.UserID = userID.(uint)
	task.Title = taskData.Title
	task.Description = taskData.Description
	task.Priority = taskData.Priority
	task.EstimatedTime = taskData.EstimatedTime
	task.Completed = taskData.Completed

	// 解析日期字符串
	if taskData.Date != "" {
		if parsedDate, err := time.Parse("2006-01-02", taskData.Date); err == nil {
			task.Date = parsedDate
		} else {
			task.Date = time.Now()
		}
	} else {
		task.Date = time.Now()
	}

	if err := pc.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// UpdateTask 更新任务
func (pc *PlannerController) UpdateTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var task models.Task
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var taskData struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Priority      int    `json:"priority"`
		EstimatedTime int    `json:"estimatedTime"`
		Date          string `json:"date"`
		Completed     bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&taskData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Title = taskData.Title
	task.Description = taskData.Description
	task.Priority = taskData.Priority
	task.EstimatedTime = taskData.EstimatedTime
	task.Completed = taskData.Completed

	// 解析日期字符串
	if taskData.Date != "" {
		if parsedDate, err := time.Parse("2006-01-02", taskData.Date); err == nil {
			task.Date = parsedDate
		}
	}

	if err := pc.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// DeleteTask 删除任务
func (pc *PlannerController) DeleteTask(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ==================== 娱乐活动相关 ====================

// GetEntertainments 获取娱乐活动列表
func (pc *PlannerController) GetEntertainments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	var entertainments []models.Entertainment
	query := pc.DB.Where("user_id = ?", userID)

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			query = query.Where("DATE(date) = ?", date.Format("2006-01-02"))
		}
	}

	if err := query.Order("created_at DESC").Find(&entertainments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entertainments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entertainments})
}

// CreateEntertainment 创建娱乐活动
func (pc *PlannerController) CreateEntertainment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var entertainmentData struct {
		Type     string   `json:"type"`
		Duration int      `json:"duration"`
		Cost     *float64 `json:"cost"`   // 使用指针，允许0值
		Rating   *int     `json:"rating"` // 使用指针，允许0值
		Notes    string   `json:"notes"`
		Date     string   `json:"date"`
	}

	if err := c.ShouldBindJSON(&entertainmentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if entertainmentData.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type is required"})
		return
	}
	if entertainmentData.Duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be greater than 0"})
		return
	}

	var entertainment models.Entertainment
	entertainment.UserID = userID.(uint)
	entertainment.Type = entertainmentData.Type
	entertainment.Duration = entertainmentData.Duration
	if entertainmentData.Cost != nil {
		entertainment.Cost = *entertainmentData.Cost
	} else {
		entertainment.Cost = 0
	}
	if entertainmentData.Rating != nil {
		entertainment.Rating = *entertainmentData.Rating
	} else {
		entertainment.Rating = 0
	}
	entertainment.Notes = entertainmentData.Notes

	// 解析日期字符串
	if entertainmentData.Date != "" {
		if parsedDate, err := time.Parse("2006-01-02", entertainmentData.Date); err == nil {
			entertainment.Date = parsedDate
		} else {
			entertainment.Date = time.Now()
		}
	} else {
		entertainment.Date = time.Now()
	}

	if err := pc.DB.Create(&entertainment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entertainment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entertainment})
}

// UpdateEntertainment 更新娱乐活动
func (pc *PlannerController) UpdateEntertainment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var entertainment models.Entertainment
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&entertainment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entertainment not found"})
		return
	}

	var entertainmentData struct {
		Type     string   `json:"type"`
		Duration int      `json:"duration"`
		Cost     *float64 `json:"cost"`
		Rating   *int     `json:"rating"`
		Notes    string   `json:"notes"`
		Date     string   `json:"date"`
	}

	if err := c.ShouldBindJSON(&entertainmentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entertainment.Type = entertainmentData.Type
	entertainment.Duration = entertainmentData.Duration
	if entertainmentData.Cost != nil {
		entertainment.Cost = *entertainmentData.Cost
	}
	if entertainmentData.Rating != nil {
		entertainment.Rating = *entertainmentData.Rating
	}
	entertainment.Notes = entertainmentData.Notes

	// 解析日期字符串
	if entertainmentData.Date != "" {
		if parsedDate, err := time.Parse("2006-01-02", entertainmentData.Date); err == nil {
			entertainment.Date = parsedDate
		}
	}

	if err := pc.DB.Save(&entertainment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entertainment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entertainment})
}

// DeleteEntertainment 删除娱乐活动
func (pc *PlannerController) DeleteEntertainment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Entertainment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entertainment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entertainment deleted successfully"})
}

// ==================== 每日反思相关 ====================

// GetReflection 获取每日反思
func (pc *PlannerController) GetReflection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	dateStr := c.Query("date")
	var date time.Time
	var err error

	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	} else {
		date = time.Now()
	}

	var reflection models.DailyReflection
	if err := pc.DB.Where("user_id = ? AND DATE(date) = ?", userID, date.Format("2006-01-02")).First(&reflection).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reflection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reflection})
}

// CreateOrUpdateReflection 创建或更新每日反思
func (pc *PlannerController) CreateOrUpdateReflection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var reflectionData struct {
		Content string             `json:"content"`
		Tags    models.StringArray `json:"tags"`
		Mood    string             `json:"mood"`
		Date    string             `json:"date"`
	}

	if err := c.ShouldBindJSON(&reflectionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reflection models.DailyReflection
	reflection.UserID = userID.(uint)
	reflection.Content = reflectionData.Content
	reflection.Tags = reflectionData.Tags
	reflection.Mood = reflectionData.Mood

	// 解析日期字符串
	if reflectionData.Date != "" {
		if parsedDate, err := time.Parse("2006-01-02", reflectionData.Date); err == nil {
			reflection.Date = parsedDate
		} else {
			reflection.Date = time.Now()
		}
	} else {
		reflection.Date = time.Now()
	}

	// 检查是否已存在
	var existing models.DailyReflection
	result := pc.DB.Where("user_id = ? AND DATE(date) = ?", userID, reflection.Date.Format("2006-01-02")).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		if err := pc.DB.Create(&reflection).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reflection"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": reflection})
	} else if result.Error == nil {
		// 更新现有记录
		existing.Content = reflection.Content
		existing.Tags = reflection.Tags
		existing.Mood = reflection.Mood
		if err := pc.DB.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reflection"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": existing})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process reflection"})
	}
}

// GetReflections 获取反思列表（用于历史记录）
func (pc *PlannerController) GetReflections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "30")
	limit := 30
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}
	var reflections []models.DailyReflection
	if err := pc.DB.Where("user_id = ?", userID).Order("date DESC").Limit(limit).Find(&reflections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reflections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reflections})
}

// ==================== 截止日期相关 ====================

// GetDeadlines 获取截止日期列表
func (pc *PlannerController) GetDeadlines(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	completed := c.Query("completed")
	var deadlines []models.Deadline
	query := pc.DB.Where("user_id = ?", userID)

	if completed == "false" {
		query = query.Where("completed = ?", false)
	}

	if err := query.Order("due_date ASC, priority DESC").Find(&deadlines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deadlines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deadlines})
}

// CreateDeadline 创建截止日期
func (pc *PlannerController) CreateDeadline(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var deadlineData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"dueDate"`
		DueTime     string `json:"dueTime"`
		Priority    int    `json:"priority"`
		Category    string `json:"category"`
		Completed   bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&deadlineData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var deadline models.Deadline
	deadline.UserID = userID.(uint)
	deadline.Title = deadlineData.Title
	deadline.Description = deadlineData.Description
	deadline.DueTime = deadlineData.DueTime
	deadline.Priority = deadlineData.Priority
	deadline.Category = deadlineData.Category
	deadline.Completed = deadlineData.Completed

	// 解析日期字符串
	if deadlineData.DueDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", deadlineData.DueDate); err == nil {
			deadline.DueDate = parsedDate
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due date is required"})
		return
	}

	if err := pc.DB.Create(&deadline).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deadline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deadline})
}

// UpdateDeadline 更新截止日期
func (pc *PlannerController) UpdateDeadline(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var deadline models.Deadline
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).First(&deadline).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deadline not found"})
		return
	}

	var deadlineData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"dueDate"`
		DueTime     string `json:"dueTime"`
		Priority    int    `json:"priority"`
		Category    string `json:"category"`
		Completed   bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&deadlineData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deadline.Title = deadlineData.Title
	deadline.Description = deadlineData.Description
	deadline.DueTime = deadlineData.DueTime
	deadline.Priority = deadlineData.Priority
	deadline.Category = deadlineData.Category
	deadline.Completed = deadlineData.Completed

	// 解析日期字符串
	if deadlineData.DueDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", deadlineData.DueDate); err == nil {
			deadline.DueDate = parsedDate
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	if err := pc.DB.Save(&deadline).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deadline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deadline})
}

// DeleteDeadline 删除截止日期
func (pc *PlannerController) DeleteDeadline(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	if err := pc.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Deadline{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deadline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deadline deleted successfully"})
}
