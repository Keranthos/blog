package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxDownloadBytes int64 = 5 * 1024 * 1024 // 5MB 上限
	httpTimeout            = 12 * time.Second
)

var allowedImageMIMEs = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// IsExternalImageURL 判断是否为需要本地化的外链图片
func IsExternalImageURL(raw string) bool {
	if raw == "" {
		return false
	}
	if strings.HasPrefix(raw, "/uploads/") || strings.HasPrefix(raw, "uploads/") {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}

// FetchAndStoreImage 下载并保存外链图片到 uploads/images，返回站内可访问路径（/uploads/images/<file>）
func FetchAndStoreImage(raw string) (string, error) {
	return FetchAndStoreImageTo(raw, "")
}

// FetchAndStoreImageTo 下载并保存外链图片到 uploads/images/<subdir>，subdir 为空则为 uploads/images 根目录
func FetchAndStoreImageTo(raw string, subdir string) (string, error) {
	if !IsExternalImageURL(raw) {
		return raw, nil
	}

	// SSRF 基础防护：禁止内网/本机地址
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return "", errors.New("blocked private address")
		}
	}

	client := &http.Client{Timeout: httpTimeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return errors.New("too many redirects")
		}
		return nil
	}}
	req, _ := http.NewRequest(http.MethodGet, raw, nil)
	// 设置更真实的User-Agent和Referer以绕过防盗链
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.weibo.com/")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// 限制体积并读取少量用于探测 MIME
	limited := io.LimitReader(resp.Body, maxDownloadBytes+1)
	buf, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	if int64(len(buf)) > maxDownloadBytes {
		return "", errors.New("image too large")
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(buf)
	}
	// 处理带参数的 content-type，比如 image/jpeg; charset=binary
	if mt, _, err := mime.ParseMediaType(mimeType); err == nil {
		mimeType = mt
	}
	ext, ok := allowedImageMIMEs[mimeType]
	if !ok {
		return "", fmt.Errorf("unsupported mime: %s", mimeType)
	}

	// 生成唯一文件名（时间戳+哈希+随机），允许同一图片多次上传
	now := time.Now()
	h := sha256.Sum256(append([]byte(raw), buf...))
	hashStr := hex.EncodeToString(h[:16]) // 使用前16字节的哈希作为部分文件名
	// 添加时间戳和纳秒级随机数确保唯一性
	filename := fmt.Sprintf("%s_%s_%d%s", now.Format("20060102_150405"), hashStr, now.Nanosecond(), ext)

	// 确保目录存在
	// 创建日期子目录（YYYY/MM）以便管理
	year := now.Format("2006")
	month := now.Format("01")

	dir := filepath.Join("uploads", "images")
	if subdir != "" {
		// 如果有子目录，添加日期子目录：uploads/images/{subdir}/YYYY/MM
		dir = filepath.Join(dir, subdir, year, month)
	} else {
		// 如果没有子目录，直接使用日期：uploads/images/YYYY/MM
		dir = filepath.Join(dir, year, month)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	abs := filepath.Join(dir, filename)

	// 直接写入文件，允许重复上传
	if err := os.WriteFile(abs, buf, 0644); err != nil {
		return "", err
	}

	// 返回站内可访问路径（与现有静态文件服务保持一致）
	return "/" + filepath.ToSlash(abs), nil
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	// 私网段
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	for _, cidr := range privateBlocks {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// StoreLocalImageFromPath 从本地路径读取图片并保存到 uploads/images，返回站内URL
func StoreLocalImageFromPath(p string) (string, error) {
	return StoreLocalImageFromPathTo(p, "")
}

// StoreLocalImageFromPathTo 从本地路径读取图片并保存到 uploads/images/<subdir>，返回站内URL
func StoreLocalImageFromPathTo(p string, subdir string) (string, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	// 猜测MIME
	mimeType := http.DetectContentType(data)
	if mt, _, err := mime.ParseMediaType(mimeType); err == nil {
		mimeType = mt
	}
	ext, ok := allowedImageMIMEs[mimeType]
	if !ok {
		// 若无匹配，尝试从扩展名推断
		if e := strings.ToLower(filepath.Ext(p)); e != "" {
			ext = e
		} else {
			return "", fmt.Errorf("unsupported mime: %s", mimeType)
		}
	}
	// 生成唯一文件名（时间戳+哈希），允许同一图片多次上传
	now := time.Now()
	h := sha256.Sum256(data)
	hashStr := hex.EncodeToString(h[:16]) // 使用前16字节的哈希作为部分文件名
	// 添加时间戳和纳秒级随机数确保唯一性
	filename := fmt.Sprintf("%s_%s_%d%s", now.Format("20060102_150405"), hashStr, now.Nanosecond(), ext)

	dir := filepath.Join("uploads", "images")
	if subdir != "" {
		dir = filepath.Join(dir, subdir)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	abs := filepath.Join(dir, filename)

	// 直接写入文件，允许重复上传
	if err := os.WriteFile(abs, data, 0644); err != nil {
		return "", err
	}
	return "/" + filepath.ToSlash(abs), nil
}
