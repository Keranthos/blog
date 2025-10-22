package models

import (
	"time"
)

// Weather 天气信息模型
type Weather struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	City        string    `json:"city" gorm:"not null"`         // 城市名称
	Country     string    `json:"country"`                      // 国家
	Latitude    float64   `json:"latitude"`                     // 纬度
	Longitude   float64   `json:"longitude"`                    // 经度
	Temperature float64   `json:"temperature"`                  // 温度
	Weather     string    `json:"weather"`                      // 天气描述
	Description string    `json:"description"`                  // 详细描述
	WindSpeed   float64   `json:"wind_speed"`                   // 风速
	Humidity    int       `json:"humidity"`                     // 湿度
	Pressure    int       `json:"pressure"`                     // 气压
	Visibility  int       `json:"visibility"`                   // 能见度
	AirQuality  string    `json:"air_quality"`                  // 空气质量
	LastUpdated time.Time `json:"last_updated" gorm:"not null"` // 最后更新时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// 新增字段
	TomorrowForecast string `json:"tomorrow_forecast"` // 明天预报
	LifeIndex        string `json:"life_index"`        // 生活指数
}

// IsToday 检查天气数据是否是今天的
func (w *Weather) IsToday() bool {
	now := time.Now()
	return w.LastUpdated.Year() == now.Year() &&
		w.LastUpdated.YearDay() == now.YearDay()
}

// IsExpired 检查天气数据是否过期（超过6小时）
func (w *Weather) IsExpired() bool {
	return time.Since(w.LastUpdated) > 6*time.Hour
}
