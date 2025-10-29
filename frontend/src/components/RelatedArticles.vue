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

// 获取类型图标
const getTypeIcon = (type) => {
  const iconMap = {
    blog: 'blog',
    project: 'code',
    research: 'flask',
    moment: 'heart'
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

  // 移除 Markdown 语法标记
  return text
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
    // 移除图片 ![alt](url) -> alt
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '$1')
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
    .trim()
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
            return {
              id: item.ID,
              type: 'moment',
              title: item.Title,
              content: item.Content,
              image: item.Image,
              tags: [],
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
            return {
              id: item.ID,
              type,
              title: item.title,
              content: rawContent,
              image: item.image,
              tags: item.tags || [],
              time: item.CreatedAt,
              viewCount: item.viewCount || 0,
              excerpt: cleanContent ? cleanContent.substring(0, 100) + (cleanContent.length > 100 ? '...' : '') : ''
            }
          })
        }

        allArticles.push(...articles)
      } catch (error) {
        console.error(`获取 ${type} 文章失败:`, error)
      }
    }

    // 过滤掉当前文章
    const filteredArticles = allArticles.filter(
      article => article.id !== props.currentArticle.id
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
  router.push(`/${article.type}/${article.id}`)
}

onMounted(() => {
  getRelatedArticles()
})
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
