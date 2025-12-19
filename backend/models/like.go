package models

import "gorm.io/gorm"

// Like 点赞模型
type Like struct {
	gorm.Model
	UserID     uint   `json:"userID" gorm:"not null;index:idx_like_user_id;uniqueIndex:idx_user_article_unique"`
	ArticleType string `json:"articleType" gorm:"type:varchar(50);not null;index:idx_like_article_type;uniqueIndex:idx_user_article_unique"`
	ArticleID  uint   `json:"articleID" gorm:"not null;index:idx_like_article_id;uniqueIndex:idx_user_article_unique"`
	User       User   `json:"user" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

