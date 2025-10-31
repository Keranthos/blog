package models

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	Poster string  `gorm:"not null"`
	Name   string  `gorm:"not null"`
	Review string  `gorm:"type:text"`
	Rating float64 `gorm:"default:8"`
	Date   string  `gorm:"not null"`
	Type   string  `gorm:"not null;index:idx_media_type"` // 添加 type 索引用于过滤查询
	// CreatedAt 索引用于排序查询
}
