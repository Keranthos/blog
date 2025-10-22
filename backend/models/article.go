package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// StringArray 自定义类型用于处理JSON字符串数组
type StringArray []string

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan into StringArray")
	}

	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口，用于向数据库写入数据
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

type BlogArticle struct {
	gorm.Model             // 自动包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string      `json:"title" gorm:"type:varchar(255);index:idx_title"`
	Content    string      `json:"content" gorm:"type:mediumtext"` // 支持最大 16MB
	Tags       StringArray `json:"tags" gorm:"type:json"`
	Image      string      `json:"image" gorm:"type:varchar(500)"`
	ViewCount  int         `json:"viewCount" gorm:"default:0;index:idx_view_count"` // 阅读量
	IsTop      bool        `json:"isTop" gorm:"default:false;index:idx_is_top"`     // 是否置顶
}

type ResearchArticle struct {
	gorm.Model             // 自动包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string      `json:"title" gorm:"type:varchar(255);index:idx_research_title"`
	Abstract   string      `json:"abstract"`
	Tags       StringArray `json:"tags" gorm:"type:json"`
	Image      string      `json:"image"`
	ViewCount  int         `json:"viewCount" gorm:"default:0;index:idx_research_view_count"` // 阅读量
	IsTop      bool        `json:"isTop" gorm:"default:false;index:idx_research_is_top"`     // 是否置顶
}

type ProjectArticle struct {
	gorm.Model             // 自动包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Title      string      `json:"title" gorm:"type:varchar(255);index:idx_project_title"`
	Status     string      `json:"status" gorm:"index:idx_project_status"`
	Tags       StringArray `json:"tags" gorm:"type:json"`
	Image      string      `json:"image"`
	ViewCount  int         `json:"viewCount" gorm:"default:0;index:idx_project_view_count"` // 阅读量
	IsTop      bool        `json:"isTop" gorm:"default:false;index:idx_project_is_top"`     // 是否置顶
}
