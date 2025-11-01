package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ExtractImageURLsFromContent 从HTML/Markdown内容中提取所有图片URL
// 支持格式：<img src="...">, ![alt](url), ![](url), markdown图片语法
func ExtractImageURLsFromContent(content string) []string {
	var imageURLs []string
	seen := make(map[string]bool)

	// 匹配HTML img标签: <img src="url"> 或 <img src='url'>
	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches := imgRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				imageURLs = append(imageURLs, url)
				seen[url] = true
			}
		}
	}

	// 匹配Markdown图片语法: ![alt](url)
	markdownRegex := regexp.MustCompile(`!\[[^\]]*\]\(([^\)]+)\)`)
	matches = markdownRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			url := strings.TrimSpace(match[1])
			if url != "" && !seen[url] {
				imageURLs = append(imageURLs, url)
				seen[url] = true
			}
		}
	}

	return imageURLs
}

// DeleteImageFileSafely 安全删除图片文件（仅删除本地文件，不检查是否被使用）
// 用于在删除文章/媒体时清理相关图片
func DeleteImageFileSafely(imagePath string) error {
	// 只处理本地路径（/uploads/images/ 开头）
	if !strings.HasPrefix(imagePath, "/uploads/images/") {
		return nil // 外链图片不需要删除
	}

	// 转换为文件系统路径
	filePath := strings.TrimPrefix(imagePath, "/")
	filePath = filepath.FromSlash(filePath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，忽略错误
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return err
	}

	// 尝试删除空的父目录
	dir := filepath.Dir(filePath)
	for strings.HasPrefix(dir, "uploads"+string(filepath.Separator)+"images") {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}
		os.Remove(dir)
		dir = filepath.Dir(dir)
	}

	return nil
}
