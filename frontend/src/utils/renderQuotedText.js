/**
 * 引用文本渲染工具
 * 用于渲染评论中引用的原文，支持Markdown和LaTeX公式
 */

import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { protectLatex, restoreAndRenderLatex } from './latex'

// 渲染缓存，避免重复渲染相同文本
const renderCache = new Map()

/**
 * 渲染引用文本为HTML
 * @param {string} quotedText - 引用的原文（Markdown格式）
 * @param {number} maxLength - 最大长度限制（字符数，保留用于未来可能的优化，当前完全依赖CSS处理截断和省略号）
 * @returns {Promise<string>} - 渲染后的HTML
 */
export async function renderQuotedText (quotedText, maxLength = 0) {
  if (!quotedText || quotedText.trim() === '') {
    return ''
  }

  // 注意：maxLength 参数保留用于未来可能的优化，但当前完全依赖CSS处理截断和省略号
  // 这样可以确保省略号正确显示，特别是对于HTML内容

  // 检查缓存
  const cacheKey = `${quotedText}_${maxLength}`
  if (renderCache.has(cacheKey)) {
    return renderCache.get(cacheKey)
  }

  try {
    // 预处理：保护 LaTeX 公式（在 marked 渲染前处理）
    const { protected: protectedContent, placeholders: latexPlaceholders } = protectLatex(quotedText)

    // 渲染 Markdown
    let html = marked(protectedContent, {
      breaks: true, // 将换行符转换为 <br>
      gfm: true, // 启用 GitHub Flavored Markdown
      headerIds: false,
      mangle: false,
      sanitize: false
    })

    // 恢复并渲染 LaTeX 公式（异步）
    html = await restoreAndRenderLatex(html, latexPlaceholders)

    // 配置DOMPurify允许style属性，确保内联样式不被过滤，支持 LaTeX 公式
    const sanitized = DOMPurify.sanitize(html, {
      ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'span', 'div', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'img', 'math', 'annotation', 'semantics', 'mrow', 'mi', 'mo', 'mn', 'mfrac', 'msup', 'msub', 'munderover', 'munder', 'mover', 'mtable', 'mtr', 'mtd', 'mtext'],
      ALLOWED_ATTR: ['class', 'style', 'src', 'alt', 'title', 'width', 'height', 'data-*', 'aria-*', 'role', 'dir', 'colspan', 'rowspan'],
      ALLOW_DATA_ATTR: true
    })

    // 缓存结果（限制缓存大小，避免内存泄漏）
    if (renderCache.size > 100) {
      // 清除最旧的缓存（简单策略：清除前50个）
      const keysToDelete = Array.from(renderCache.keys()).slice(0, 50)
      keysToDelete.forEach(key => renderCache.delete(key))
    }
    renderCache.set(cacheKey, sanitized)

    return sanitized
  } catch (error) {
    console.error('渲染引用文本失败:', error)
    // 出错时返回纯文本
    return quotedText
  }
}
