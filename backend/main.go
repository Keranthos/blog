package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"
)

func main() {
	//初始化配置
	config.InitConfig()

	if config.AppConfig == nil {
		log.Fatalf("bull shit!AppConfig is not initialized")
	}

	// 初始化uploads目录
	initUploadsDir()

	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务器
	port := config.AppConfig.App.Port
	if port == "" {
		port = "8080" // 默认端口
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// initUploadsDir 初始化uploads目录
func initUploadsDir() {
	// 创建uploads目录结构
	dirs := []string{
		"uploads",
		"uploads/images",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create directory %s: %v", dir, err)
		} else {
			log.Printf("Created directory: %s", dir)
		}
	}
}
