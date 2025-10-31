package models

import "gorm.io/gorm"

type Moment struct {
	gorm.Model
	Title     string `gorm:"type:varchar(255);not null;index:idx_moment_title" json:"Title"`
	Content   string `gorm:"type:text;not null" json:"Content"`
	Image     string `gorm:"type:varchar(255)" json:"Image"`
	Tags      string `gorm:"type:varchar(500)" json:"Tags"`
	Author    string `gorm:"type:varchar(100);index:idx_moment_author" json:"Author"`
	ViewCount int    `json:"viewCount" gorm:"default:0;index:idx_moment_view_count"` // 阅读量
	// CreatedAt 索引用于排序查询
}
