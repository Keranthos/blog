package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// 验证用户名
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}

	if utf8.RuneCountInString(username) > 12 {
		return errors.New("用户名最多12个字符")
	}

	if utf8.RuneCountInString(username) < 2 {
		return errors.New("用户名至少2个字符")
	}

	// 检查是否包含非法字符
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\u4e00-\u9fa5]+$`, username)
	if !matched {
		return errors.New("用户名只能包含字母、数字、下划线和中文")
	}

	return nil
}

// 验证密码
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	if len(password) < 6 {
		return errors.New("密码至少6位")
	}

	if len(password) > 50 {
		return errors.New("密码最多50位")
	}

	// 检查密码强度
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	strength := 0
	if hasLower {
		strength++
	}
	if hasUpper {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength++
	}

	if strength < 2 {
		return errors.New("密码强度不够，建议包含大小写字母、数字和特殊字符")
	}

	return nil
}

// 验证邮箱
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// 验证文章标题
func ValidateArticleTitle(title string) error {
	if title == "" {
		return errors.New("标题不能为空")
	}

	if utf8.RuneCountInString(title) > 100 {
		return errors.New("标题最多100个字符")
	}

	if utf8.RuneCountInString(title) < 2 {
		return errors.New("标题至少2个字符")
	}

	// 移除首尾空格
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("标题不能只包含空格")
	}

	return nil
}

// 验证文章内容
func ValidateArticleContent(content string) error {
	if content == "" {
		return errors.New("内容不能为空")
	}

	if utf8.RuneCountInString(content) > 50000 {
		return errors.New("内容最多50000个字符")
	}

	if utf8.RuneCountInString(content) < 10 {
		return errors.New("内容至少10个字符")
	}

	// 移除首尾空格
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("内容不能只包含空格")
	}

	return nil
}

// 验证标签
func ValidateTags(tags []string) error {
	if len(tags) > 10 {
		return errors.New("标签最多10个")
	}

	for _, tag := range tags {
		if utf8.RuneCountInString(tag) > 20 {
			return errors.New("单个标签最多20个字符")
		}

		if utf8.RuneCountInString(tag) < 1 {
			return errors.New("标签不能为空")
		}

		// 检查标签是否包含非法字符
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\u4e00-\u9fa5]+$`, tag)
		if !matched {
			return errors.New("标签只能包含字母、数字、下划线和中文")
		}
	}

	return nil
}

// 验证评论内容
func ValidateCommentContent(content string) error {
	if content == "" {
		return errors.New("评论内容不能为空")
	}

	if utf8.RuneCountInString(content) > 1000 {
		return errors.New("评论内容最多1000个字符")
	}

	if utf8.RuneCountInString(content) < 2 {
		return errors.New("评论内容至少2个字符")
	}

	// 移除首尾空格
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("评论内容不能只包含空格")
	}

	return nil
}

// 验证问题内容
func ValidateQuestionContent(content string) error {
	if content == "" {
		return errors.New("问题内容不能为空")
	}

	if utf8.RuneCountInString(content) > 500 {
		return errors.New("问题内容最多500个字符")
	}

	if utf8.RuneCountInString(content) < 5 {
		return errors.New("问题内容至少5个字符")
	}

	// 移除首尾空格
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("问题内容不能只包含空格")
	}

	return nil
}

// 验证回答内容
func ValidateAnswerContent(content string) error {
	if content == "" {
		return errors.New("回答内容不能为空")
	}

	if utf8.RuneCountInString(content) > 2000 {
		return errors.New("回答内容最多2000个字符")
	}

	if utf8.RuneCountInString(content) < 5 {
		return errors.New("回答内容至少5个字符")
	}

	// 移除首尾空格
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("回答内容不能只包含空格")
	}

	return nil
}

// 验证图片URL
func ValidateImageURL(url string) error {
	if url == "" {
		return nil // 图片URL可以为空
	}

	// 检查URL格式
	urlRegex := regexp.MustCompile(`^https?:\/\/[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		return errors.New("图片URL格式不正确")
	}

	if len(url) > 500 {
		return errors.New("图片URL最多500个字符")
	}

	return nil
}

// 清理和标准化输入
func SanitizeInput(input string) string {
	// 移除首尾空格
	input = strings.TrimSpace(input)

	// 移除多余的换行符和空格
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	return input
}

// 检查是否包含恶意内容
func ContainsMaliciousContent(content string) bool {
	// 简单的恶意内容检测（可以扩展）
	maliciousPatterns := []string{
		"<script",
		"javascript:",
		"onload=",
		"onerror=",
		"<iframe",
		"<object",
		"<embed",
	}

	content = strings.ToLower(content)
	for _, pattern := range maliciousPatterns {
		if strings.Contains(content, pattern) {
			return true
		}
	}

	return false
}
