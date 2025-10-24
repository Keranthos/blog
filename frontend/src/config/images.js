/**
 * 图片配置文件
 * 使用 Unsplash API 提供高质量图片
 *
 * Unsplash 使用说明：
 * https://source.unsplash.com/{尺寸}/?{关键词}
 *
 * 常用关键词：
 * - technology, coding, computer (技术类)
 * - science, laboratory, research (科研类)
 * - nature, landscape, mountain (自然风景)
 * - city, architecture, building (城市建筑)
 * - abstract, minimal, gradient (抽象简约)
 * - books, library, reading (书籍阅读)
 * - workspace, desk, office (工作空间)
 */

export const imageConfig = {
  // 页面头图配置（固定图片，不会变化）
  headers: {
    blog: 'https://picsum.photos/id/180/1920/1080', // 科技感的键盘
    project: 'https://picsum.photos/id/201/1920/1080', // 工作空间/办公桌
    research: 'https://picsum.photos/id/225/1920/1080', // 实验室/科学设备
    questionBox: 'https://picsum.photos/id/1043/1920/1080', // 问号/创意/灯光
    books: 'https://picsum.photos/id/24/1920/1080', // 书籍/图书馆
    novels: 'https://picsum.photos/id/159/1920/1080', // 阅读/书本
    movies: 'https://picsum.photos/id/146/1920/1080', // 电影/胶片
    moments: 'https://picsum.photos/id/201/1920/1080' // 碎碎念（与项目不同的图）
  },

  // 文章卡片默认图片（当作者未指定时使用）
  defaultArticleImage: 'https://picsum.photos/id/1/800/600',

  // 精选文章封面图片（固定ID，适合不同主题）
  presetArticleImages: [
    // 科技类
    { id: 180, name: '科技键盘', url: 'https://picsum.photos/id/180/800/600', category: 'tech' },
    { id: 119, name: '电脑代码', url: 'https://picsum.photos/id/119/800/600', category: 'tech' },
    { id: 0, name: '笔记本电脑', url: 'https://picsum.photos/id/0/800/600', category: 'tech' },

    // 工作空间
    { id: 201, name: '工作空间', url: 'https://picsum.photos/id/201/800/600', category: 'workspace' },
    { id: 7, name: '办公桌', url: 'https://picsum.photos/id/7/800/600', category: 'workspace' },
    { id: 30, name: '咖啡桌', url: 'https://picsum.photos/id/30/800/600', category: 'workspace' },

    // 科研学术
    { id: 225, name: '实验室', url: 'https://picsum.photos/id/225/800/600', category: 'research' },
    { id: 24, name: '书籍', url: 'https://picsum.photos/id/24/800/600', category: 'research' },
    { id: 159, name: '笔记本', url: 'https://picsum.photos/id/159/800/600', category: 'research' },

    // 自然风景
    { id: 10, name: '自然风景', url: 'https://picsum.photos/id/10/800/600', category: 'nature' },
    { id: 15, name: '山景', url: 'https://picsum.photos/id/15/800/600', category: 'nature' },
    { id: 40, name: '海洋', url: 'https://picsum.photos/id/40/800/600', category: 'nature' },
    { id: 20, name: '森林', url: 'https://picsum.photos/id/20/800/600', category: 'nature' },

    // 城市建筑
    { id: 1, name: '城市建筑', url: 'https://picsum.photos/id/1/800/600', category: 'urban' },
    { id: 8, name: '夜景', url: 'https://picsum.photos/id/8/800/600', category: 'urban' },
    { id: 3, name: '建筑物', url: 'https://picsum.photos/id/3/800/600', category: 'urban' },

    // 抽象艺术
    { id: 60, name: '抽象艺术', url: 'https://picsum.photos/id/60/800/600', category: 'abstract' },
    { id: 152, name: '创意灯泡', url: 'https://picsum.photos/id/152/800/600', category: 'abstract' },
    { id: 103, name: '艺术图案', url: 'https://picsum.photos/id/103/800/600', category: 'abstract' }
  ],

  // 备用图片源
  fallback: {
    // Picsum Photos (Lorem Picsum)
    picsum: {
      random: (width = 1920, height = 1080) => `https://picsum.photos/${width}/${height}`,
      blur: (width = 1920, height = 1080) => `https://picsum.photos/${width}/${height}?blur=2`,
      grayscale: (width = 1920, height = 1080) => `https://picsum.photos/${width}/${height}?grayscale`
    },

    // 纯色渐变背景（CSS）
    gradients: {
      purple: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      blue: 'linear-gradient(135deg, #667eea 0%, #4facfe 100%)',
      sunset: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
      ocean: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
      forest: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
    }
  },

  // 不同主题的推荐关键词
  themes: {
    tech: ['technology', 'coding', 'computer', 'digital', 'programming'],
    nature: ['nature', 'landscape', 'mountain', 'forest', 'ocean'],
    minimal: ['abstract', 'minimal', 'gradient', 'simple', 'clean'],
    urban: ['city', 'architecture', 'building', 'street', 'urban'],
    creative: ['art', 'creative', 'design', 'colorful', 'artistic']
  }
}

/**
 * 生成随机 Unsplash 图片 URL
 * @param {number} width - 图片宽度
 * @param {number} height - 图片高度
 * @param {string|string[]} keywords - 关键词（字符串或数组）
 * @returns {string} 图片 URL
 */
export function getUnsplashImage (width = 1920, height = 1080, keywords = 'nature') {
  const keywordString = Array.isArray(keywords) ? keywords.join(',') : keywords
  return `https://source.unsplash.com/${width}x${height}/?${keywordString}`
}

/**
 * 生成 Picsum 随机图片 URL
 * @param {number} width - 图片宽度
 * @param {number} height - 图片高度
 * @param {Object} options - 选项 {blur: boolean, grayscale: boolean}
 * @returns {string} 图片 URL
 */
export function getPicsumImage (width = 1920, height = 1080, options = {}) {
  let url = `https://picsum.photos/${width}/${height}`
  const params = []

  if (options.blur) params.push('blur=2')
  if (options.grayscale) params.push('grayscale')

  if (params.length > 0) {
    url += '?' + params.join('&')
  }

  return url
}

export default imageConfig
