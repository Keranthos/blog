package models

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Presentation 讲演模型
type Presentation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:255;not null" validate:"required"`
	Tags      string    `json:"tags" gorm:"type:text"`              // JSON格式存储标签数组
	FilePath  string    `json:"file_path" gorm:"size:500;not null"` // PDF文件路径
	Thumbnail string    `json:"thumbnail" gorm:"size:500"`          // 缩略图路径
	FileSize  int64     `json:"file_size"`                          // 文件大小(字节)
	FileType  string    `json:"file_type" gorm:"size:50"`           // 文件类型(pdf)
	ViewCount int       `json:"view_count" gorm:"default:0"`        // 浏览次数
	IsActive  bool      `json:"is_active" gorm:"default:true"`      // 是否激活
	IsPrivate bool      `json:"is_private" gorm:"default:false"`    // 是否私密
	Password  string    `json:"password" gorm:"size:255"`           // 访问密码
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id" gorm:"not null"`       // 创建者ID
	User      User      `json:"user" gorm:"foreignKey:UserID"` // 关联用户
}

// PresentationCreateRequest 创建讲演请求
type PresentationCreateRequest struct {
	Title     string `json:"title" validate:"required"`
	Tags      string `json:"tags"`       // JSON格式的标签数组
	IsPrivate bool   `json:"is_private"` // 是否私密
}

// PresentationUpdateRequest 更新讲演请求
type PresentationUpdateRequest struct {
	Title     string `json:"title"`
	Tags      string `json:"tags"`
	IsActive  *bool  `json:"is_active"`
	IsPrivate *bool  `json:"is_private"` // 是否私密
}

// PresentationResponse 讲演响应
type PresentationResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	FilePath  string    `json:"file_path"`
	Thumbnail string    `json:"thumbnail"`
	FileSize  int64     `json:"file_size"`
	FileType  string    `json:"file_type"`
	ViewCount int       `json:"view_count"`
	IsActive  bool      `json:"is_active"`
	IsPrivate bool      `json:"is_private"` // 是否私密
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
}

// ToResponse 转换为响应格式
func (p *Presentation) ToResponse() PresentationResponse {
	var tags []string
	if p.Tags != "" {
		// 解析JSON字符串为数组
		if err := json.Unmarshal([]byte(p.Tags), &tags); err != nil {
			// 如果JSON解析失败，尝试按逗号分割
			if strings.Contains(p.Tags, ",") {
				tags = strings.Split(p.Tags, ",")
				// 去除每个标签的前后空格
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			} else {
				// 如果既不是JSON也不是逗号分割，就当作单个标签
				tags = []string{strings.TrimSpace(p.Tags)}
			}
		}
	}

	// 处理缩略图URL
	var thumbnailURL string
	if p.Thumbnail != "" {
		// 如果缩略图是URL（以http开头），直接使用
		if strings.HasPrefix(p.Thumbnail, "http") {
			thumbnailURL = p.Thumbnail
		} else {
			// 如果是本地文件路径，构建完整的URL
			filename := filepath.Base(p.Thumbnail)
			thumbnailURL = fmt.Sprintf("http://localhost:3000/uploads/presentations/%s", filename)
		}
	}

	return PresentationResponse{
		ID:        p.ID,
		Title:     p.Title,
		Tags:      tags,
		FilePath:  p.FilePath,
		Thumbnail: thumbnailURL,
		FileSize:  p.FileSize,
		FileType:  p.FileType,
		ViewCount: p.ViewCount,
		IsActive:  p.IsActive,
		IsPrivate: p.IsPrivate,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		User:      p.User,
	}
}
