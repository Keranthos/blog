package config

import (
	"log"
)

// CreateIndexes 创建必要的数据库索引以优化查询性能
func CreateIndexes() error {
	log.Println("开始创建数据库索引...")

	// 为 BlogArticle 创建 created_at 索引（用于排序查询）
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_blog_created_at ON blog_articles(created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_blog_created_at 索引失败: %v", err)
	}

	// 为 ResearchArticle 创建 created_at 索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_research_created_at ON research_articles(created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_research_created_at 索引失败: %v", err)
	}

	// 为 ProjectArticle 创建 created_at 索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_project_created_at ON project_articles(created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_project_created_at 索引失败: %v", err)
	}

	// 为 Moment 创建 created_at 索引
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_moment_created_at ON moments(created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_moment_created_at 索引失败: %v", err)
	}

	// 为 Media 创建 created_at 索引（用于排序查询）
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_media_created_at ON media(created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_media_created_at 索引失败: %v", err)
	}

	// 为 Media 创建复合索引 (type, created_at) 用于按类型筛选和排序
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_media_type_created_at ON media(type, created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_media_type_created_at 索引失败: %v", err)
	}

	// 为 Comment 创建复合索引 (blog_id, type, created_at) 用于按文章筛选和排序
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_comment_blog_type_created_at ON comments(blog_id, type, created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_comment_blog_type_created_at 索引失败: %v", err)
	}

	// 为 Comment 创建复合索引 (parent_id, created_at) 用于回复查询
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_comment_parent_created_at ON comments(parent_id, created_at DESC)").Error; err != nil {
		log.Printf("警告: 创建 idx_comment_parent_created_at 索引失败: %v", err)
	}

	log.Println("数据库索引创建完成")
	return nil
}
