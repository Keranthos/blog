<template>
  <div v-if="relatedArticles.length > 0" class="related-articles">
    <div class="related-header">
      <h3>
        <font-awesome-icon icon="link" class="header-icon" />
        相关文章
      </h3>
      <p class="related-subtitle">你可能还对这些文章感兴趣</p>
    </div>

    <div class="related-list">
      <div
        v-for="article in relatedArticles"
        :key="article.id"
        class="related-item"
        @click="goToArticle(article)"
      >
        <div class="related-image">
          <img
            :src="article.image || getDefaultImage(article.type)"
            :alt="article.title"
            loading="lazy"
            decoding="async"
            @error="onImgError($event)"
          />
        </div>

        <div class="related-content">
          <h4 class="related-title">
            <font-awesome-icon :icon="getTypeIcon(article.type)" class="title-icon" />
            {{ article.title }}
          </h4>
          <p v-if="article.excerpt" class="related-excerpt">{{ article.excerpt }}</p>

          <div class="related-meta">
            <div class="meta-left">
              <span class="meta-item">
                <font-awesome-icon icon="calendar" />
                {{ formatDate(article.time) }}
              </span>
              <span v-if="article.viewCount" class="meta-item">
                <font-awesome-icon icon="eye" />
                {{ article.viewCount }}
              </span>
            </div>
            <div v-if="article.tags && article.tags.length" class="related-tags">
              <span
                v-for="tag in article.tags.slice(0, 3)"
                :key="tag"
                class="tag"
              >
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticlesList } from '@/api/Articles/browse'
import { getMomentsList } from '@/api/Moments/browse'

const props = defineProps({
  currentArticle: {
    type: Object,
    required: true
  },
  maxCount: {
    type: Number,
    default: 3
  }
})

const router = useRouter()
const relatedArticles = ref([])
const loading = ref(false)
const fallbackImg = '/images/sunset-mountains.jpg'

// 获取类型图标
const getTypeIcon = (type) => {
  const iconMap = {
    blog: 'blog',
    project: 'code',
    research: 'flask',
    moment: 'comment-dots'
  }
  return iconMap[type] || 'file'
}

// 获取默认图片
const getDefaultImage = (type) => {
  const defaultImages = {
    blog: 'https://picsum.photos/id/1/400/300',
    project: 'https://picsum.photos/id/2/400/300',
    research: 'https://picsum.photos/id/3/400/300',
    moment: 'https://picsum.photos/id/4/400/300'
  }
  return defaultImages[type] || 'https://picsum.photos/id/5/400/300'
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 去除 Markdown 语法标记
const stripMarkdown = (text) => {
  if (!text) return ''

  // 调试：记录原始文本
  if (text.includes('paste_1761985840582.png') || text.includes('!paste')) {
    console.log('🔍 [stripMarkdown] 原始文本:', text.substring(0, 200))
  }

  // 先移除图片相关内容（在其他Markdown处理之前，避免下划线被误删）
  // 1. 完全移除图片语法 ![alt](url) -> 空
  let cleaned = text.replace(/!\[([^\]]*)\]\([^)]+\)/g, '')

  // 2. 移除完整的HTTP/HTTPS图片URL
  cleaned = cleaned.replace(/https?:\/\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')

  // 3. 移除相对路径图片（/uploads/, /images/）
  cleaned = cleaned.replace(/\/(uploads|images)\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')

  // 4. 移除带 ! 前缀的图片文件名（如 !paste_1761985840582.png）
  cleaned = cleaned.replace(/!\s*[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')

  // 5. 移除纯图片文件名（使用单词边界，避免误删）
  cleaned = cleaned.replace(/\b[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')

  // 调试：检查图片移除后
  if (text.includes('paste_1761985840582.png') || text.includes('!paste')) {
    console.log('🔍 [stripMarkdown] 图片移除后:', cleaned.substring(0, 200))
  }

  // 现在处理其他 Markdown 语法标记
  cleaned = cleaned
    // 移除标题标记 (##, ###等)
    .replace(/^#{1,6}\s+/gm, '')
    // 移除粗体 (**text** 或 __text__)
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/__([^_]+)__/g, '$1')
    // 移除斜体 (*text* 或 _text_)
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/_([^_]+)_/g, '$1')
    // 移除删除线 (~~text~~)
    .replace(/~~([^~]+)~~/g, '$1')
    // 移除行内代码 (`code`)
    .replace(/`([^`]+)`/g, '$1')
    // 移除链接 [text](url) -> text
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    // 移除引用 (> text)
    .replace(/^>\s+/gm, '')
    // 移除列表标记 (-, *, +, 数字.)
    .replace(/^[\s]*[-*+]\s+/gm, '')
    .replace(/^[\s]*\d+\.\s+/gm, '')
    // 移除代码块标记 (```language 或 ```)
    .replace(/```[\w]*\n[\s\S]*?```/g, '')
    // 移除水平线 (---, ***)
    .replace(/^[-*]{3,}$/gm, '')
    // 移除多余的空行
    .replace(/\n{3,}/g, '\n\n')

  // 调试：测试正则表达式匹配
  if (text.includes('paste_1761985840582.png') || text.includes('!paste') || text.includes('paste1761985840582')) {
    console.log('🔍 [stripMarkdown] 处理前文本片段:', cleaned.substring(0, 300))
    // 检查是否还有残留
    const hasImageUrl = /(https?:\/\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(\/(uploads|images)\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(!\s*[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(\b[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))/gi.test(cleaned)
    console.log('🔍 [stripMarkdown] 是否还包含图片URL:', hasImageUrl)
    if (hasImageUrl) {
      const remainingMatches = cleaned.match(/(https?:\/\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(\/(uploads|images)\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(!\s*[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))|(\b[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico))/gi)
      console.log('🔍 [stripMarkdown] 残留的图片URL:', remainingMatches)
    }
  }

  // 调试：检查最终结果
  if (text.includes('paste_1761985840582.png') || text.includes('!paste')) {
    console.log('🔍 [stripMarkdown] 处理后文本片段:', cleaned.substring(0, 300))
    console.log('🔍 [stripMarkdown] 是否还包含图片URL:', cleaned.includes('paste_1761985840582.png') || cleaned.includes('!paste'))
  }

  // 清理移除URL后可能出现的多余空格（保留换行）
  cleaned = cleaned.replace(/[ \t]+/g, ' ').replace(/\n{3,}/g, '\n\n')

  return cleaned.trim()
}

// 计算文章相似度
const calculateSimilarity = (article1, article2) => {
  let score = 0

  // 标签相似度 (权重最高)
  if (article1.tags && article2.tags) {
    const commonTags = article1.tags.filter(tag =>
      article2.tags.some(t => t.toLowerCase() === tag.toLowerCase())
    )
    score += commonTags.length * 10
  }

  // 标题相似度
  if (article1.title && article2.title) {
    const title1 = article1.title.toLowerCase()
    const title2 = article2.title.toLowerCase()

    // 检查共同词汇
    const words1 = title1.split(/\s+/)
    const words2 = title2.split(/\s+/)
    const commonWords = words1.filter(word =>
      words2.some(w => w.includes(word) || word.includes(w))
    )
    score += commonWords.length * 3
  }

  // 内容相似度 (基于关键词)
  if (article1.content && article2.content) {
    const content1 = article1.content.toLowerCase()
    const content2 = article2.content.toLowerCase()

    // 简单关键词匹配
    const keywords = ['技术', '开发', '学习', '经验', '项目', '问题', '解决', '方法']
    keywords.forEach(keyword => {
      if (content1.includes(keyword) && content2.includes(keyword)) {
        score += 1
      }
    })
  }

  // 类型相同加分
  if (article1.type === article2.type) {
    score += 2
  }

  return score
}

// 获取相关文章
const getRelatedArticles = async () => {
  if (!props.currentArticle) return

  loading.value = true

  try {
    // 获取所有文章
    const allArticles = []
    const types = ['blog', 'project', 'research', 'moment']

    for (const type of types) {
      try {
        let articles = []
        if (type === 'moment') {
          const response = await getMomentsList(1, 50)
          articles = response.data.map(item => {
            const cleanContent = stripMarkdown(item.Content || '')
            // 解析标签（Tags可能是逗号分隔的字符串）
            let tags = []
            if (item.Tags && typeof item.Tags === 'string') {
              tags = item.Tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0)
            } else if (Array.isArray(item.Tags)) {
              tags = item.Tags
            }
            return {
              id: String(item.ID),
              type: 'moment',
              title: item.Title,
              content: item.Content,
              image: item.Image,
              tags: tags,
              time: item.CreatedAt,
              viewCount: item.ViewCount || 0,
              excerpt: cleanContent ? cleanContent.substring(0, 100) + (cleanContent.length > 100 ? '...' : '') : ''
            }
          })
        } else {
          const response = await getArticlesList(type, 1, 50)
          articles = response.data.map(item => {
            const rawContent = item.content || item.abstract || ''
            const cleanContent = stripMarkdown(rawContent)
            let excerpt = cleanContent ? cleanContent.substring(0, 100) + (cleanContent.length > 100 ? '...' : '') : ''

            // 再次清理excerpt，确保移除任何残留的图片URL
            excerpt = excerpt
              .replace(/!\s*[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')
              .replace(/https?:\/\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')
              .replace(/\/(uploads|images)\/[^\s)]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')
              .replace(/\b[a-zA-Z0-9_.-]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico)(\?[^\s)]*)?/gi, '')
              .replace(/\s+/g, ' ')
              .trim()

            // 调试：检查生成的excerpt
            if (rawContent.includes('paste_1761985840582.png') || rawContent.includes('!paste') || excerpt.includes('paste') || excerpt.includes('.png')) {
              console.log('🔍 [getRelatedArticles] 文章ID:', item.ID)
              console.log('🔍 [getRelatedArticles] 原始内容片段:', rawContent.substring(0, 150))
              console.log('🔍 [getRelatedArticles] 清理后内容片段:', cleanContent.substring(0, 150))
              console.log('🔍 [getRelatedArticles] 截取后的excerpt:', cleanContent ? cleanContent.substring(0, 100) + (cleanContent.length > 100 ? '...' : '') : '')
              console.log('🔍 [getRelatedArticles] 最终excerpt:', excerpt)
              console.log('🔍 [getRelatedArticles] excerpt是否包含图片URL:', excerpt.includes('paste') || excerpt.includes('.png'))
            }

            return {
              id: String(item.ID),
              type,
              title: item.title,
              content: rawContent,
              image: item.image,
              tags: item.tags || [],
              time: item.CreatedAt,
              viewCount: item.viewCount || 0,
              excerpt
            }
          })
        }

        allArticles.push(...articles)
      } catch (error) {
        console.error(`获取 ${type} 文章失败:`, error)
      }
    }

    // 过滤掉当前文章（使用字符串比较避免类型不一致）
    const currentArticleId = String(props.currentArticle.id)
    const filteredArticles = allArticles.filter(
      article => String(article.id) !== currentArticleId
    )

    // 计算相似度并排序
    const articlesWithScore = filteredArticles.map(article => ({
      ...article,
      similarity: calculateSimilarity(props.currentArticle, article)
    }))

    // 按相似度排序，取前几个
    relatedArticles.value = articlesWithScore
      .sort((a, b) => b.similarity - a.similarity)
      .slice(0, props.maxCount)
      .filter(article => article.similarity > 0) // 只显示有相似度的文章
  } catch (error) {
    console.error('获取相关文章失败:', error)
  } finally {
    loading.value = false
  }
}

// 跳转到文章
const goToArticle = (article) => {
  if (article.type === 'moment') {
    router.push(`/moments/${article.id}`)
  } else if (article.type === 'blog') {
    router.push({ name: 'BlogDetail', params: { id: article.id } })
  } else if (article.type === 'research' || article.type === 'project') {
    router.push({ name: 'BlogDetail', params: { id: article.id }, query: { type: article.type } })
  }
}

onMounted(() => {
  getRelatedArticles()
})

// 图片错误回退
const onImgError = (e) => {
  const img = e?.target
  if (img && img.src !== fallbackImg) {
    img.src = fallbackImg
  }
}
</script>

<style scoped>
.related-articles {
  margin-top: 30px;
  padding: 30px 0;
  background: transparent;
  border-radius: 0;
  box-shadow: none;
  border: none;
  backdrop-filter: none;
  margin-bottom: 30px;
}

.related-header {
  text-align: left;
  margin-bottom: 30px;
}

.related-header h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
}

.header-icon {
  color: #667eea;
  font-size: 1.3rem;
}

.related-subtitle {
  color: #666;
  font-size: 0.95rem;
  margin: 0;
}

.related-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.related-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  background: transparent;
  border-radius: 0;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  backdrop-filter: none;
  margin-bottom: 0;
  box-shadow: none;
}

.related-item:last-child {
  border-bottom: none;
}

.related-item:hover {
  background: rgba(255, 255, 255, 0.05);
  transform: none;
  box-shadow: none;
  border-bottom-color: rgba(102, 126, 234, 0.2);
  backdrop-filter: blur(10px);
}

.related-image {
  position: relative;
  flex-shrink: 0;
  width: 120px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
}

.related-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.related-item:hover .related-image img {
  transform: scale(1.1);
}

.related-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.related-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin: 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 8px;
  text-align: left;
}

.title-icon {
  color: #667eea;
  font-size: 1rem;
  flex-shrink: 0;
}

.related-excerpt {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-align: left;
}

.related-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 15px;
  color: #999;
  font-size: 0.85rem;
}

.meta-left {
  display: flex;
  gap: 15px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.related-tags {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
}

.tag {
  padding: 2px 8px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .related-articles {
    margin-top: 30px;
    padding: 20px;
  }

  .related-item {
    flex-direction: column;
    gap: 15px;
  }

  .related-image {
    width: 100%;
    height: 120px;
  }

  .related-header h3 {
    font-size: 1.3rem;
  }
}

@media (max-width: 480px) {
  .related-articles {
    padding: 15px;
  }

  .related-item {
    padding: 15px;
  }

  .related-image {
    height: 100px;
  }

  .related-title {
    font-size: 1rem;
  }

  .related-meta {
    gap: 10px;
  }
}
</style>
