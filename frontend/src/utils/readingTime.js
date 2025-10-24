/**
 * 文章阅读时间估算工具
 * 基于中英文阅读速度差异进行估算
 */

// 中文阅读速度：每分钟300-500字
const CHINESE_WPM = 400 // 中文每分钟字数
// 英文阅读速度：每分钟200-300字
const ENGLISH_WPM = 250 // 英文每分钟字数

/**
 * 估算文章阅读时间
 * @param {string} content - 文章内容
 * @param {string} title - 文章标题（可选）
 * @returns {Object} 包含阅读时间的对象
 */
export function estimateReadingTime (content, title = '') {
  if (!content) {
    return {
      minutes: 0,
      words: 0,
      display: '小于 1 分钟'
    }
  }

  // 清理HTML标签，获取纯文本
  const cleanContent = stripHtmlTags(content + ' ' + title)

  // 统计中英文字符
  const stats = analyzeText(cleanContent)

  // 计算阅读时间（分钟）
  const chineseTime = stats.chineseChars / CHINESE_WPM
  const englishTime = stats.englishWords / ENGLISH_WPM
  const totalMinutes = chineseTime + englishTime

  // 至少1分钟
  const finalMinutes = Math.max(1, Math.ceil(totalMinutes))

  return {
    minutes: finalMinutes,
    words: stats.totalWords,
    chineseChars: stats.chineseChars,
    englishWords: stats.englishWords,
    display: formatReadingTime(finalMinutes)
  }
}

/**
 * 移除HTML标签
 * @param {string} html - 包含HTML标签的文本
 * @returns {string} 纯文本
 */
function stripHtmlTags (html) {
  // 移除HTML标签
  let text = html.replace(/<[^>]*>/g, '')

  // 移除多余的空白字符
  text = text.replace(/\s+/g, ' ').trim()

  // 移除代码块中的特殊字符
  text = text.replace(/```[\s\S]*?```/g, '')
  text = text.replace(/`[^`]*`/g, '')

  return text
}

/**
 * 分析文本内容
 * @param {string} text - 纯文本内容
 * @returns {Object} 文本统计信息
 */
function analyzeText (text) {
  let chineseChars = 0
  let englishWords = 0

  // 中文字符正则：[\u4e00-\u9fff]
  const chineseRegex = /[\u4e00-\u9fff]/g
  const chineseMatches = text.match(chineseRegex)
  chineseChars = chineseMatches ? chineseMatches.length : 0

  // 英文字符和数字
  const englishText = text.replace(/[\u4e00-\u9fff]/g, ' ')
  const englishWordsArray = englishText.match(/\b[a-zA-Z0-9]+\b/g)
  englishWords = englishWordsArray ? englishWordsArray.length : 0

  // 其他字符（标点符号等）
  const otherChars = text.length - chineseChars - englishWords

  return {
    chineseChars,
    englishWords,
    otherChars,
    totalWords: chineseChars + englishWords + Math.floor(otherChars / 3) // 标点符号折算
  }
}

/**
 * 格式化阅读时间显示
 * @param {number} minutes - 分钟数
 * @returns {string} 格式化的时间字符串
 */
function formatReadingTime (minutes) {
  if (minutes < 1) {
    return '小于 1 分钟'
  } else if (minutes === 1) {
    return '1 分钟'
  } else if (minutes < 60) {
    return `${minutes} 分钟`
  } else {
    const hours = Math.floor(minutes / 60)
    const remainingMinutes = minutes % 60

    if (remainingMinutes === 0) {
      return `${hours} 小时`
    } else {
      return `${hours} 小时 ${remainingMinutes} 分钟`
    }
  }
}

/**
 * 获取阅读时间图标
 * @param {number} minutes - 分钟数
 * @returns {string} 图标名称
 */
export function getReadingTimeIcon (minutes) {
  if (minutes < 5) {
    return 'bolt' // 闪电 - 快速阅读
  } else if (minutes < 15) {
    return 'clock' // 时钟 - 中等阅读
  } else {
    return 'book-open' // 书本 - 深度阅读
  }
}

/**
 * 获取阅读时间颜色
 * @param {number} minutes - 分钟数
 * @returns {string} CSS颜色值
 */
export function getReadingTimeColor (minutes) {
  if (minutes < 5) {
    return '#4CAF50' // 绿色 - 快速
  } else if (minutes < 15) {
    return '#FF9800' // 橙色 - 中等
  } else {
    return '#F44336' // 红色 - 深度
  }
}

/**
 * 批量估算多篇文章的阅读时间
 * @param {Array} articles - 文章数组
 * @returns {Array} 包含阅读时间的文章数组
 */
export function batchEstimateReadingTime (articles) {
  return articles.map(article => ({
    ...article,
    readingTime: estimateReadingTime(article.content || article.Content, article.title || article.Title)
  }))
}
