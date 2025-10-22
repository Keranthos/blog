package controllers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/models"
)

// RSS结构体定义
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Category    string `xml:"category,omitempty"`
}

// 生成RSS订阅源
func GenerateRSS(c *gin.Context) {
	// 获取所有博客文章
	var blogArticles []models.BlogArticle
	var projectArticles []models.ProjectArticle
	var researchArticles []models.ResearchArticle
	var moments []models.Moment

	// 并行获取所有类型的文章
	blogChan := make(chan []models.BlogArticle, 1)
	projectChan := make(chan []models.ProjectArticle, 1)
	researchChan := make(chan []models.ResearchArticle, 1)
	momentsChan := make(chan []models.Moment, 1)

	// 获取博客文章
	go func() {
		var articles []models.BlogArticle
		config.DB.Order("created_at DESC").Limit(20).Find(&articles)
		blogChan <- articles
	}()

	// 获取项目文章
	go func() {
		var articles []models.ProjectArticle
		config.DB.Order("created_at DESC").Limit(10).Find(&articles)
		projectChan <- articles
	}()

	// 获取科研文章
	go func() {
		var articles []models.ResearchArticle
		config.DB.Order("created_at DESC").Limit(10).Find(&articles)
		researchChan <- articles
	}()

	// 获取碎碎念
	go func() {
		var articles []models.Moment
		config.DB.Order("created_at DESC").Limit(10).Find(&articles)
		momentsChan <- articles
	}()

	// 等待所有goroutine完成
	blogArticles = <-blogChan
	projectArticles = <-projectChan
	researchArticles = <-researchChan
	moments = <-momentsChan

	// 合并所有文章
	var allItems []RSSItem

	// 处理博客文章
	for _, article := range blogArticles {
		item := RSSItem{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/blog/%d", getBaseURL(c), article.ID),
			Description: truncateText(article.Content, 200),
			PubDate:     formatRSSDate(article.CreatedAt),
			GUID:        fmt.Sprintf("blog-%d", article.ID),
			Category:    "博客",
		}
		allItems = append(allItems, item)
	}

	// 处理项目文章
	for _, article := range projectArticles {
		item := RSSItem{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/project/%d", getBaseURL(c), article.ID),
			Description: fmt.Sprintf("项目状态: %s", article.Status),
			PubDate:     formatRSSDate(article.CreatedAt),
			GUID:        fmt.Sprintf("project-%d", article.ID),
			Category:    "项目",
		}
		allItems = append(allItems, item)
	}

	// 处理科研文章
	for _, article := range researchArticles {
		item := RSSItem{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/research/%d", getBaseURL(c), article.ID),
			Description: truncateText(article.Abstract, 200),
			PubDate:     formatRSSDate(article.CreatedAt),
			GUID:        fmt.Sprintf("research-%d", article.ID),
			Category:    "科研",
		}
		allItems = append(allItems, item)
	}

	// 处理碎碎念
	for _, moment := range moments {
		item := RSSItem{
			Title:       moment.Title,
			Link:        fmt.Sprintf("%s/moments/%d", getBaseURL(c), moment.ID),
			Description: truncateText(moment.Content, 200),
			PubDate:     formatRSSDate(moment.CreatedAt),
			GUID:        fmt.Sprintf("moment-%d", moment.ID),
			Category:    "碎碎念",
		}
		allItems = append(allItems, item)
	}

	// 按时间排序（最新的在前）
	sortRSSItemsByDate(allItems)

	// 限制RSS项目数量
	if len(allItems) > 50 {
		allItems = allItems[:50]
	}

	// 构建RSS结构
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         "山角函兽的个人博客",
			Link:          getBaseURL(c),
			Description:   "技术博客、项目分享、科研记录、生活感悟",
			Language:      "zh-cn",
			LastBuildDate: formatRSSDate(time.Now()),
			Items:         allItems,
		},
	}

	// 设置响应头
	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600") // 缓存1小时

	// 输出XML
	c.XML(http.StatusOK, rss)
}

// 获取基础URL
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	return fmt.Sprintf("%s://%s", scheme, host)
}

// 格式化RSS日期
func formatRSSDate(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// 截断文本
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// 按日期排序RSS项目
func sortRSSItemsByDate(items []RSSItem) {
	// 简单的冒泡排序（按PubDate降序）
	for i := 0; i < len(items)-1; i++ {
		for j := 0; j < len(items)-i-1; j++ {
			date1, err1 := time.Parse(time.RFC1123Z, items[j].PubDate)
			date2, err2 := time.Parse(time.RFC1123Z, items[j+1].PubDate)

			if err1 != nil || err2 != nil {
				continue
			}

			if date1.Before(date2) {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}
}

// 获取RSS统计信息
func GetRSSStats(c *gin.Context) {
	var blogCount, projectCount, researchCount, momentCount int64

	// 并行获取统计数据
	blogChan := make(chan int64, 1)
	projectChan := make(chan int64, 1)
	researchChan := make(chan int64, 1)
	momentsChan := make(chan int64, 1)

	go func() {
		var count int64
		config.DB.Model(&models.BlogArticle{}).Count(&count)
		blogChan <- count
	}()

	go func() {
		var count int64
		config.DB.Model(&models.ProjectArticle{}).Count(&count)
		projectChan <- count
	}()

	go func() {
		var count int64
		config.DB.Model(&models.ResearchArticle{}).Count(&count)
		researchChan <- count
	}()

	go func() {
		var count int64
		config.DB.Model(&models.Moment{}).Count(&count)
		momentsChan <- count
	}()

	blogCount = <-blogChan
	projectCount = <-projectChan
	researchCount = <-researchChan
	momentCount = <-momentsChan

	totalCount := blogCount + projectCount + researchCount + momentCount

	c.JSON(http.StatusOK, gin.H{
		"total":    totalCount,
		"blog":     blogCount,
		"project":  projectCount,
		"research": researchCount,
		"moment":   momentCount,
	})
}

// 生成sitemap.xml
func GenerateSitemap(c *gin.Context) {
	baseURL := getBaseURL(c)

	// 获取所有文章数据
	var blogArticles []models.BlogArticle
	var projectArticles []models.ProjectArticle
	var researchArticles []models.ResearchArticle
	var moments []models.Moment

	// 并行获取数据
	blogChan := make(chan []models.BlogArticle, 1)
	projectChan := make(chan []models.ProjectArticle, 1)
	researchChan := make(chan []models.ResearchArticle, 1)
	momentsChan := make(chan []models.Moment, 1)

	go func() {
		var articles []models.BlogArticle
		config.DB.Order("updated_at DESC").Find(&articles)
		blogChan <- articles
	}()

	go func() {
		var articles []models.ProjectArticle
		config.DB.Order("updated_at DESC").Find(&articles)
		projectChan <- articles
	}()

	go func() {
		var articles []models.ResearchArticle
		config.DB.Order("updated_at DESC").Find(&articles)
		researchChan <- articles
	}()

	go func() {
		var articles []models.Moment
		config.DB.Order("updated_at DESC").Find(&articles)
		momentsChan <- articles
	}()

	blogArticles = <-blogChan
	projectArticles = <-projectChan
	researchArticles = <-researchChan
	moments = <-momentsChan

	// 生成sitemap XML
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	// 添加静态页面
	staticPages := []struct {
		URL        string
		Priority   string
		ChangeFreq string
	}{
		{"/", "1.0", "daily"},
		{"/blog", "0.8", "weekly"},
		{"/project", "0.8", "weekly"},
		{"/research", "0.8", "weekly"},
		{"/moments", "0.8", "weekly"},
		{"/tags", "0.6", "weekly"},
		{"/timeline", "0.6", "weekly"},
		{"/archive", "0.6", "weekly"},
		{"/questionbox", "0.5", "weekly"},
	}

	for _, page := range staticPages {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>%s</changefreq>
    <priority>%s</priority>
  </url>`, baseURL, page.URL, time.Now().Format("2006-01-02"), page.ChangeFreq, page.Priority)
	}

	// 添加博客文章
	for _, article := range blogArticles {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/blog/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>`, baseURL, article.ID, article.UpdatedAt.Format("2006-01-02"))
	}

	// 添加项目文章
	for _, article := range projectArticles {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/project/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>`, baseURL, article.ID, article.UpdatedAt.Format("2006-01-02"))
	}

	// 添加科研文章
	for _, article := range researchArticles {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/research/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>`, baseURL, article.ID, article.UpdatedAt.Format("2006-01-02"))
	}

	// 添加碎碎念
	for _, moment := range moments {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/moments/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>`, baseURL, moment.ID, moment.UpdatedAt.Format("2006-01-02"))
	}

	sitemap += `
</urlset>`

	// 设置响应头
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600") // 缓存1小时
	c.String(http.StatusOK, sitemap)
}
