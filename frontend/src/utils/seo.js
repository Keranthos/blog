/**
 * SEO工具类 - 用于动态更新页面meta标签
 */

// 默认SEO信息
const defaultSEO = {
  title: '山角函兽的个人博客 - 记录生活、分享技术、探索未知',
  description: '山角函兽的个人博客 - 记录生活、分享技术、探索未知。包含技术博客、项目分享、科研记录、生活感悟等内容。',
  keywords: '个人博客,技术博客,项目分享,科研记录,生活感悟,山角函兽',
  image: '',
  url: ''
}

/**
 * 更新页面SEO信息
 * @param {Object} seoData - SEO数据对象
 */
export function updateSEO (seoData) {
  const data = { ...defaultSEO, ...seoData }

  // 更新标题
  if (data.title) {
    document.title = data.title
    updateMetaTag('property', 'og:title', data.title)
    updateMetaTag('name', 'twitter:title', data.title)
  }

  // 更新描述
  if (data.description) {
    updateMetaTag('name', 'description', data.description)
    updateMetaTag('property', 'og:description', data.description)
    updateMetaTag('name', 'twitter:description', data.description)
  }

  // 更新关键词
  if (data.keywords) {
    updateMetaTag('name', 'keywords', data.keywords)
  }

  // 更新图片
  if (data.image) {
    updateMetaTag('property', 'og:image', data.image)
    updateMetaTag('name', 'twitter:image', data.image)
  }

  // 更新URL
  if (data.url) {
    updateMetaTag('property', 'og:url', data.url)
    updateMetaTag('rel', 'canonical', data.url)
  }

  // 更新文章相关meta
  if (data.type) {
    updateMetaTag('property', 'og:type', data.type)
  }

  if (data.publishedTime) {
    updateMetaTag('property', 'article:published_time', data.publishedTime)
  }

  if (data.modifiedTime) {
    updateMetaTag('property', 'article:modified_time', data.modifiedTime)
  }

  if (data.author) {
    updateMetaTag('property', 'article:author', data.author)
  }

  if (data.section) {
    updateMetaTag('property', 'article:section', data.section)
  }

  if (data.tags && Array.isArray(data.tags)) {
    data.tags.forEach(tag => {
      updateMetaTag('property', 'article:tag', tag)
    })
  }
}

/**
 * 更新或创建meta标签
 * @param {string} attribute - 属性名 (name 或 property)
 * @param {string} value - 属性值
 * @param {string} content - 内容
 */
function updateMetaTag (attribute, value, content) {
  let meta = document.querySelector(`meta[${attribute}="${value}"]`)

  if (!meta) {
    meta = document.createElement('meta')
    meta.setAttribute(attribute, value)
    document.head.appendChild(meta)
  }

  meta.setAttribute('content', content)
}

/**
 * 为文章页面生成SEO信息
 * @param {Object} article - 文章对象
 * @param {string} type - 文章类型
 */
export function generateArticleSEO (article, type) {
  const typeMap = {
    blog: '博客',
    project: '项目',
    research: '科研',
    moment: '碎碎念'
  }

  const typeName = typeMap[type] || type
  const title = `${article.title} - 山角函兽的${typeName}`

  // 生成描述（从内容中提取或使用摘要）
  let description = article.abstract || article.content || article.status || ''
  if (description.length > 150) {
    description = description.substring(0, 150) + '...'
  }
  description = `${description} - 来自山角函兽的${typeName}文章`

  // 生成关键词
  let keywords = `山角函兽,${typeName}`
  if (article.tags && Array.isArray(article.tags)) {
    keywords += ',' + article.tags.join(',')
  }

  // 生成图片URL
  let image = article.image || ''
  if (image && !image.startsWith('http')) {
    image = window.location.origin + image
  }

  // 生成文章URL
  const url = `${window.location.origin}/${type}/${article.ID}`

  // 格式化时间
  const publishedTime = new Date(article.CreatedAt).toISOString()
  const modifiedTime = new Date(article.UpdatedAt).toISOString()

  return {
    title,
    description,
    keywords,
    image,
    url,
    type: 'article',
    publishedTime,
    modifiedTime,
    author: '山角函兽',
    section: typeName,
    tags: article.tags || []
  }
}

/**
 * 为列表页面生成SEO信息
 * @param {string} pageType - 页面类型
 * @param {number} page - 页码
 */
export function generateListSEO (pageType, page = 1) {
  const typeMap = {
    blog: { name: '博客', desc: '技术博客文章' },
    project: { name: '项目', desc: '项目分享' },
    research: { name: '科研', desc: '科研记录' },
    moment: { name: '碎碎念', desc: '生活感悟' },
    tags: { name: '标签云', desc: '文章标签分类' },
    timeline: { name: '时间树', desc: '按时间浏览文章' },
    archive: { name: '文章归档', desc: '按年份月份浏览文章' },
    search: { name: '搜索', desc: '搜索文章内容' }
  }

  const typeInfo = typeMap[pageType] || { name: pageType, desc: '' }
  const title = page > 1
    ? `${typeInfo.name} - 第${page}页 - 山角函兽的博客`
    : `${typeInfo.name} - 山角函兽的博客`

  const description = `${typeInfo.desc} - 山角函兽的个人博客，记录生活、分享技术、探索未知`

  return {
    title,
    description,
    keywords: `山角函兽,${typeInfo.name},个人博客`,
    url: `${window.location.origin}/${pageType}${page > 1 ? `?page=${page}` : ''}`
  }
}

/**
 * 重置为默认SEO
 */
export function resetSEO () {
  updateSEO(defaultSEO)
}
