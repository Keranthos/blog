package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	BlogID     uint      `gorm:"not null;index:idx_comment_blog_id" json:"blogID"`
	Username   string    `gorm:"not null;index:idx_comment_username" json:"username"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Type       string    `gorm:"not null;index:idx_comment_type" json:"type"`
	ParentID   *uint     `json:"parent_id" gorm:"default:null;index:idx_comment_parent_id"` // 父评论ID，支持嵌套回复
	QuotedText *string   `gorm:"type:text" json:"quoted_text"`                              // 引用的原文（可选）
	Replies    []Comment `json:"replies" gorm:"foreignKey:ParentID"`                        // 子评论
}
