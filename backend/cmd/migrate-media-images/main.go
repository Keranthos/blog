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
			fmt.Printf("跳过 ID=%d: 封面为空\n", m.ID)
			continue
		}
		if !utils.IsExternalImageURL(m.Poster) {
			// 已是站内URL
			skipped++
			fmt.Printf("跳过 ID=%d: 已是本地图片 (%s)\n", m.ID, m.Poster)
			continue
		}
		fmt.Printf("正在迁移 ID=%d: %s\n", m.ID, m.Poster)
		local, err := utils.FetchAndStoreImageTo(m.Poster, "media")
		if err != nil {
			failed++
			fmt.Printf("失败 ID=%d: %v\n", m.ID, err)
			continue
		}
		if err := config.DB.Model(&models.Media{}).Where("id = ?", m.ID).Update("poster", local).Error; err != nil {
			failed++
			fmt.Printf("更新数据库失败 ID=%d: %v\n", m.ID, err)
			continue
		}
		migrated++
		fmt.Printf("成功 ID=%d: %s -> %s\n", m.ID, m.Poster, local)
	}

	fmt.Printf("Media migration done. total=%d migrated=%d skipped=%d failed=%d\n", len(items), migrated, skipped, failed)
}
