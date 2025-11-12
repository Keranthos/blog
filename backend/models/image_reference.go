package models

import "gorm.io/gorm"

// ImageReference 记录图片与内容之间的引用关系
type ImageReference struct {
	gorm.Model
	ImagePath string `gorm:"type:varchar(500);index:idx_image_path"`
	RefType   string `gorm:"type:varchar(32);index:idx_ref_type_id,priority:1"`
	RefID     uint   `gorm:"index:idx_ref_type_id,priority:2"`
	UsageType string `gorm:"type:varchar(32);index"`
}
