package main

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
)

func main() {
	// 初始化配置与数据库
	config.InitConfig()
	if err := config.InitDB(); err != nil {
		panic(fmt.Errorf("init db failed: %w", err))
	}

	var items []models.Media
	if err := config.DB.Find(&items).Error; err != nil {
		panic(fmt.Errorf("load media failed: %w", err))
	}

	migrated := 0
	skipped := 0
	failed := 0

	for i := range items {
		m := &items[i]
		if m.Poster == "" {
			skipped++
			continue
		}
		local, err := utils.FetchAndStoreImageTo(m.Poster, "media")
		if err != nil {
			failed++
			continue
		}
		if local == m.Poster {
			// 已是站内URL
			skipped++
			continue
		}
		if err := config.DB.Model(&models.Media{}).Where("id = ?", m.ID).Update("poster", local).Error; err != nil {
			failed++
			continue
		}
		migrated++
	}

	fmt.Printf("Media migration done. total=%d migrated=%d skipped=%d failed=%d\n", len(items), migrated, skipped, failed)
}
