package imageref

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/models"

	"gorm.io/gorm"
)

var imagePathPattern = regexp.MustCompile(`(?i)(https?://[^\s"'()]+/uploads/[^\s"'()]+|/uploads/[^\s"'()]+|\\uploads\\[^\s"'()]+|uploads/[^\s"'()]+)`)

// NormalizeImagePath 将图片路径标准化为统一的 /uploads/... 形式，仅处理本地图片
func NormalizeImagePath(raw string) string {
	if raw == "" {
		return ""
	}

	path := strings.TrimSpace(raw)
	path = strings.ReplaceAll(path, "\\", "/")

	if strings.HasPrefix(path, "data:") {
		return ""
	}

	if strings.HasPrefix(path, "uploads/") {
		path = "/" + path
	} else {
		uploadsIdx := strings.Index(path, "/uploads/")
		if uploadsIdx >= 0 {
			path = path[uploadsIdx:]
		} else if strings.HasPrefix(path, "/uploads/") {
			// already normalized
		} else {
			return ""
		}
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}

// extractImagePaths 从内容中提取本地图片路径
func extractImagePaths(content string) []string {
	if content == "" {
		return nil
	}

	matches := imagePathPattern.FindAllString(content, -1)
	if len(matches) == 0 {
		return nil
	}

	unique := make(map[string]struct{})
	for _, match := range matches {
		path := NormalizeImagePath(match)
		if path == "" {
			continue
		}
		unique[path] = struct{}{}
	}

	paths := make([]string, 0, len(unique))
	for path := range unique {
		paths = append(paths, path)
	}
	return paths
}

// SyncImageReferences 同步指定内容的图片引用
func SyncImageReferences(db *gorm.DB, refType string, refID uint, cover string, content string, extraPaths ...string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ref_type = ? AND ref_id = ?", refType, refID).Delete(&models.ImageReference{}).Error; err != nil {
			return err
		}

		now := time.Now()
		var records []models.ImageReference

		if coverPath := NormalizeImagePath(cover); coverPath != "" {
			records = append(records, models.ImageReference{
				ImagePath: coverPath,
				RefType:   refType,
				RefID:     refID,
				UsageType: "cover",
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
			})
		}

		contentPaths := extractImagePaths(content)
		if len(extraPaths) > 0 {
			for _, ext := range extraPaths {
				if normalized := NormalizeImagePath(ext); normalized != "" {
					contentPaths = append(contentPaths, normalized)
				}
			}
		}

		if len(contentPaths) > 0 {
			unique := make(map[string]struct{})
			for _, p := range contentPaths {
				if p == "" {
					continue
				}
				if _, exists := unique[p]; exists {
					continue
				}
				unique[p] = struct{}{}
				records = append(records, models.ImageReference{
					ImagePath: p,
					RefType:   refType,
					RefID:     refID,
					UsageType: "content",
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
				})
			}
		}

		if len(records) > 0 {
			if err := tx.Create(&records).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveImageReferences 删除指定内容的所有图片引用
func RemoveImageReferences(db *gorm.DB, refType string, refID uint) error {
	return db.Where("ref_type = ? AND ref_id = ?", refType, refID).Delete(&models.ImageReference{}).Error
}

// EnsureBackfillIfEmpty 当引用表为空时，回填现有数据
func EnsureBackfillIfEmpty(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.ImageReference{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	type article struct {
		ID      uint
		Image   string
		Content string
	}

	type media struct {
		ID     uint
		Poster string
		Review string
		Type   string
	}

	var blogs []article
	if err := db.Model(&models.BlogArticle{}).Select("id, image, content").Find(&blogs).Error; err != nil {
		return err
	}
	for _, b := range blogs {
		_ = SyncImageReferences(db, "blog", b.ID, b.Image, b.Content)
	}

	var researches []article
	if err := db.Model(&models.ResearchArticle{}).Select("id, image, content").Find(&researches).Error; err != nil {
		return err
	}
	for _, r := range researches {
		_ = SyncImageReferences(db, "research", r.ID, r.Image, r.Content)
	}

	var projects []article
	if err := db.Model(&models.ProjectArticle{}).Select("id, image, content").Find(&projects).Error; err != nil {
		return err
	}
	for _, p := range projects {
		_ = SyncImageReferences(db, "project", p.ID, p.Image, p.Content)
	}

	var moments []struct {
		ID      uint
		Image   string
		Content string
	}
	if err := db.Model(&models.Moment{}).Select("id, image, content").Find(&moments).Error; err != nil {
		return err
	}
	for _, m := range moments {
		_ = SyncImageReferences(db, "moment", m.ID, m.Image, m.Content)
	}

	var mediaItems []media
	if err := db.Model(&models.Media{}).Select("id, poster, review, type").Find(&mediaItems).Error; err != nil {
		return err
	}
	for _, m := range mediaItems {
		_ = SyncImageReferences(db, "media", m.ID, m.Poster, m.Review)
	}

	return nil
}

// ListReferencesByImagePath 查询某个图片路径的所有引用记录，返回标准化后的路径
func ListReferencesByImagePath(db *gorm.DB, rawPath string) ([]models.ImageReference, string, error) {
	normalized := NormalizeImagePath(rawPath)
	if normalized == "" {
		return nil, "", nil
	}
	var refs []models.ImageReference
	if err := db.Where("image_path = ?", normalized).Find(&refs).Error; err != nil {
		return nil, "", err
	}
	return refs, normalized, nil
}

// AggregateUsageCounts 根据引用记录生成使用摘要（分类型和用途计数）
func AggregateUsageCounts(refs []models.ImageReference) map[string]map[string]int {
	result := make(map[string]map[string]int)
	for _, ref := range refs {
		if _, ok := result[ref.RefType]; !ok {
			result[ref.RefType] = make(map[string]int)
		}
		key := ref.UsageType + ":" + strconv.FormatUint(uint64(ref.RefID), 10)
		result[ref.RefType][key]++
	}
	return result
}
