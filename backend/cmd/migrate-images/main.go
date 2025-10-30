package main

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"os"
	"path/filepath"
)

func migrateTables() error {
	type item struct{ table string }
	tables := []item{{"blog_articles"}, {"research_articles"}, {"project_articles"}, {"moments"}}
	total := 0
	updated := 0
	for _, t := range tables {
		rows, err := config.DB.Table(t.table).Select("id, image").Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id uint
			var image string
			if err := rows.Scan(&id, &image); err != nil {
				continue
			}
			total++
			if !utils.IsExternalImageURL(image) {
				continue
			}
			if newURL, err := utils.FetchAndStoreImage(image); err == nil && newURL != "" {
				config.DB.Table(t.table).Where("id = ?", id).Update("image", newURL)
				updated++
			}
		}
	}
	fmt.Printf("Scanned %d records, updated %d external images.\n", total, updated)
	return nil
}

func replaceSpecificBlogCover() error {
	var blog models.BlogArticle
	if err := config.DB.Where("title = ?", "建立一个项目2_GORM入门").First(&blog).Error; err != nil {
		fmt.Println("Target blog not found or DB error:", err)
		return nil
	}
	// 根据当前工作目录定位到前端资源
	root, _ := os.Getwd()
	candidates := []string{
		filepath.Join(root, "frontend", "src", "assets", "blog_background.jpg"),
		filepath.Join(root, "..", "frontend", "src", "assets", "blog_background.jpg"),
		filepath.Join(root, "..", "..", "frontend", "src", "assets", "blog_background.jpg"),
	}
	var asset string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			asset = p
			break
		}
	}
	if asset == "" {
		return fmt.Errorf("blog header asset not found under frontend/src/assets/blog_background.jpg")
	}
	if newURL, err := utils.StoreLocalImageFromPath(asset); err == nil && newURL != "" {
		config.DB.Model(&blog).Update("image", newURL)
		fmt.Println("Updated target blog cover to:", newURL)
	} else if err != nil {
		return err
	}
	return nil
}

func main() {
	// 初始化配置与数据库
	config.InitConfig()
	if err := config.InitDB(); err != nil {
		fmt.Println("InitDB error:", err)
		os.Exit(1)
	}
	if err := migrateTables(); err != nil {
		fmt.Println("migrateTables error:", err)
		os.Exit(1)
	}
	if err := replaceSpecificBlogCover(); err != nil {
		fmt.Println("replaceSpecificBlogCover error:", err)
		os.Exit(1)
	}
	fmt.Println("Migration finished.")
}
