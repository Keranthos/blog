package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// 验证用户名
func ValidateUsername(username string) error {
	// 先去除首尾空格
	username = strings.TrimSpace(username)
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
	// 使用编译好的正则表达式提高性能
	// 使用 \p{Han} 匹配中文字符（CJK统一汉字），这是Go正则表达式推荐的方式
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\p{Han}]+$`)
	if !usernameRegex.MatchString(username) {
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

	// 标签可以为空数组
	if len(tags) == 0 {
		return nil
	}

	// 正则表达式：允许字母、数字、下划线、中文字符
	// \p{Han} 表示所有中文字符（CJK统一汉字）
	// 注意：不包含空格，标签内不能有空格（但标签本身可以有空格作为分隔）
	tagRegex := regexp.MustCompile(`^[a-zA-Z0-9_\p{Han}]+$`)

	for i, tag := range tags {
		// 先去除首尾空格
		tag = strings.TrimSpace(tag)

		if utf8.RuneCountInString(tag) < 1 {
			return fmt.Errorf("标签不能为空（索引: %d）", i)
		}

		if utf8.RuneCountInString(tag) > 20 {
			return fmt.Errorf("单个标签最多20个字符（当前标签: '%s'，长度: %d）", tag, utf8.RuneCountInString(tag))
		}

		// 检查标签是否包含非法字符
		// 允许：字母、数字、下划线、中文字符（不含空格）
		matched := tagRegex.MatchString(tag)
		if !matched {
			// 输出更详细的错误信息以便调试
			return fmt.Errorf("标签只能包含字母、数字、下划线和中文（当前标签[%d]: '%s'，长度: %d，字节: %v）",
				i, tag, utf8.RuneCountInString(tag), []byte(tag))
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

	if len(url) > 500 {
		return errors.New("图片URL最多500个字符")
	}

	// 支持三种格式：
	// 1. HTTP/HTTPS 外链: https://example.com/image.jpg
	// 2. 本地绝对路径: /uploads/images/xxx.jpg
	// 3. 本地相对路径: uploads/images/xxx.jpg
	httpRegex := regexp.MustCompile(`^https?:\/\/[^\s/$.?#].[^\s]*$`)
	localAbsoluteRegex := regexp.MustCompile(`^\/uploads\/[^\s]*$`)
	localRelativeRegex := regexp.MustCompile(`^uploads\/[^\s]*$`)

	if httpRegex.MatchString(url) || localAbsoluteRegex.MatchString(url) || localRelativeRegex.MatchString(url) {
		return nil
	}

	return errors.New("图片URL格式不正确（支持：http/https外链或/uploads/本地路径）")
}

// 清理和标准化输入
func SanitizeInput(input string) string {
	// 移除首尾空格
	input = strings.TrimSpace(input)

	// 移除多余的换行符和空格
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")

	return input
}

// 清理引用文本输入（保留换行符、代码块和LaTeX公式，移除其他Markdown格式字符）
func SanitizeQuotedText(input string) string {
	// 移除首尾空格和换行符
	input = strings.TrimSpace(input)

	// 检查是否包含代码块或LaTeX公式
	hasCodeBlock := strings.Contains(input, "```")
	hasLatex := strings.Contains(input, "$$") || strings.Contains(input, "$")

	// 如果包含代码块或LaTeX公式，保留原始格式（不清理）
	if hasCodeBlock || hasLatex {
		// 只清理多余的连续空格（保留换行符）
		input = regexp.MustCompile(`[ \t]+`).ReplaceAllString(input, " ")
		// 清理多个连续换行符（最多保留两个）
		input = regexp.MustCompile(`\n{3,}`).ReplaceAllString(input, "\n\n")
		return input
	}

	// 如果不包含代码块或LaTeX公式，移除所有Markdown格式字符，只保留纯文本和换行符
	// 移除粗体 **text** 或 __text__
	input = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(input, "$1")
	input = regexp.MustCompile(`__([^_]+)__`).ReplaceAllString(input, "$1")
	
	// 移除斜体 *text* 或 _text_
	input = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(input, "$1")
	input = regexp.MustCompile(`_([^_]+)_`).ReplaceAllString(input, "$1")
	
	// 移除链接 [text](url)
	input = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(input, "$1")
	
	// 移除图片 ![alt](url)
	input = regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`).ReplaceAllString(input, "$1")
	
	// 移除行内代码 `code`
	input = regexp.MustCompile("`([^`]+)`").ReplaceAllString(input, "$1")
	
	// 移除标题标记 # ## ### 等
	input = regexp.MustCompile(`^#{1,6}\s+`).ReplaceAllString(input, "")
	
	// 移除引用标记 >
	input = regexp.MustCompile(`^>\s+`).ReplaceAllString(input, "")
	
	// 移除列表标记 - * + 或数字.
	input = regexp.MustCompile(`^[\s]*[-*+]\s+`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`^[\s]*\d+\.\s+`).ReplaceAllString(input, "")
	
	// 移除删除线 ~~text~~
	input = regexp.MustCompile(`~~([^~]+)~~`).ReplaceAllString(input, "$1")
	
	// 只清理多余的连续空格（保留换行符）
	// 将多个连续空格替换为单个空格，但保留换行符
	input = regexp.MustCompile(`[ \t]+`).ReplaceAllString(input, " ")
	
	// 清理多个连续换行符（最多保留两个）
	input = regexp.MustCompile(`\n{3,}`).ReplaceAllString(input, "\n\n")

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
