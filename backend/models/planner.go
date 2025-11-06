package models

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	gorm.Model
	UserID        uint      `json:"userId" gorm:"not null;index:idx_user_date"`         // 用户ID
	Date          time.Time `json:"date" gorm:"type:date;not null;index:idx_user_date"` // 任务日期
	Title         string    `json:"title" gorm:"type:varchar(255);not null"`            // 任务标题
	Description   string    `json:"description" gorm:"type:text"`                       // 任务描述
	Completed     bool      `json:"completed" gorm:"default:false;index:idx_completed"` // 是否完成
	Priority      int       `json:"priority" gorm:"default:1"`                          // 优先级：1-低，2-中，3-高
	EstimatedTime int       `json:"estimatedTime" gorm:"default:0"`                     // 预计时间（分钟）
	WeekDay       int       `json:"weekDay"`                                            // 周视图中的星期几（0-6，0=周日）
}

// Entertainment 娱乐活动模型
type Entertainment struct {
	gorm.Model
	UserID   uint      `json:"userId" gorm:"not null;index:idx_user_date"`         // 用户ID
	Date     time.Time `json:"date" gorm:"type:date;not null;index:idx_user_date"` // 活动日期
	Type     string    `json:"type" gorm:"type:varchar(50);not null"`              // 活动类型（游戏、追剧、阅读等）
	Duration int       `json:"duration" gorm:"not null"`                           // 时长（分钟）
	Cost     float64   `json:"cost" gorm:"default:0"`                              // 花费（可选）
	Rating   int       `json:"rating" gorm:"default:0"`                            // 评价（1-5星，可选）
	Notes    string    `json:"notes" gorm:"type:text"`                             // 备注
}

// DailyReflection 每日反思模型
type DailyReflection struct {
	gorm.Model
	UserID  uint        `json:"userId" gorm:"not null;index:idx_user_date"`               // 用户ID
	Date    time.Time   `json:"date" gorm:"type:date;not null;uniqueIndex:idx_user_date"` // 日期（每天只能有一条）
	Content string      `json:"content" gorm:"type:mediumtext;not null"`                  // 反思内容（支持Markdown）
	Tags    StringArray `json:"tags" gorm:"type:json"`                                    // 标签（工作、学习、生活等）
	Mood    string      `json:"mood" gorm:"type:varchar(20)"`                             // 心情（可选）
}

// Deadline 截止日期模型
type Deadline struct {
	gorm.Model
	UserID      uint      `json:"userId" gorm:"not null;index:idx_user_date"` // 用户ID
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`    // 标题
	Description string    `json:"description" gorm:"type:text"`               // 描述
	DueDate     time.Time `json:"dueDate" gorm:"not null;index:idx_due_date"` // 截止日期
	DueTime     string    `json:"dueTime" gorm:"type:varchar(10)"`            // 截止时间（如：23:59）
	Completed   bool      `json:"completed" gorm:"default:false"`             // 是否完成
	Priority    int       `json:"priority" gorm:"default:1"`                  // 优先级：1-低，2-中，3-高
	Category    string    `json:"category" gorm:"type:varchar(50)"`           // 分类（工作、学习、生活等）
}
