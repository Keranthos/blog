<template>
  <div class="comments-view">
    <NavBar />
    <div class="content-container">
      <div class="comments-header">
        <h1 class="page-title">评论</h1>
        <p class="page-subtitle">最新的评论动态</p>
      </div>

      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="comments.length === 0" class="empty-container">
        <p>暂无评论</p>
      </div>

      <div v-else class="comments-list">
        <div v-for="comment in comments" :key="comment.ID" class="comment-card">
          <!-- 评论内容 -->
          <div class="comment-content">
            <div class="comment-text" v-html="renderCommentContent(comment.content)"></div>
          </div>

          <!-- 文章信息 -->
          <div class="article-info">
            <div class="article-type-icon">
              <font-awesome-icon :icon="getArticleTypeIcon(comment.articleType)" />
            </div>
            <router-link
              :to="getArticleLink(comment.articleType, comment.blogID)"
              class="article-title"
            >
              {{ comment.articleTitle || '未知文章' }}
            </router-link>
          </div>

          <!-- 评论元信息 -->
          <div class="comment-meta">
            <span class="comment-author">{{ comment.username }}</span>
            <span class="comment-time">
              <font-awesome-icon icon="clock" />
              {{ formatCommentTime(comment.CreatedAt) }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import NavBar from '@/components/NavBar.vue'
import { getAllComments } from '@/api/Comments/browse'
import { marked } from 'marked'

const comments = ref([])
const loading = ref(true)

// 加载评论列表
const loadComments = async () => {
  try {
    loading.value = true
    const data = await getAllComments()
    comments.value = data || []
  } catch (error) {
    console.error('加载评论失败:', error)
    comments.value = []
  } finally {
    loading.value = false
  }
}

// 格式化评论时间
const formatCommentTime = (timestamp) => {
  if (!timestamp) return '未知时间'
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) {
    return '刚刚'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)}小时前`
  } else if (diff < 604800000) {
    return `${Math.floor(diff / 86400000)}天前`
  } else {
    const year = date.getFullYear()
    const month = (date.getMonth() + 1).toString().padStart(2, '0')
    const day = date.getDate().toString().padStart(2, '0')
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
  }
}

// 获取文章类型图标
const getArticleTypeIcon = (type) => {
  const iconMap = {
    blog: 'blog',
    project: 'code',
    research: 'flask',
    moment: 'pen-to-square'
  }
  return iconMap[type] || 'file'
}

// 获取文章链接
const getArticleLink = (type, blogID) => {
  if (type === 'moment') {
    return `/moments/${blogID}`
  } else if (type === 'blog') {
    return `/blog/${blogID}`
  } else if (type === 'research') {
    // 如果research有独立路由，使用独立路由；否则使用通用路由
    return `/research/${blogID}`
  } else if (type === 'project') {
    // 如果project有独立路由，使用独立路由；否则使用通用路由
    return `/project/${blogID}`
  }
  return '/'
}

// 渲染评论内容（支持Markdown）
const renderCommentContent = (content) => {
  if (!content) return ''
  // 将换行符转换为HTML
  const processedContent = content.replace(/\n/g, '<br>')
  return marked(processedContent)
}

onMounted(() => {
  loadComments()
})
</script>

<style scoped>
.comments-view {
  min-height: 100vh;
  padding-top: 80px;
}

.content-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.comments-header {
  margin-bottom: 2rem;
  text-align: center;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 1rem;
  color: #666;
}

.loading-container,
.empty-container {
  text-align: center;
  padding: 4rem 2rem;
  color: #999;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(102, 126, 234, 0.1);
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.comment-card:hover {
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.15);
  transform: translateY(-2px);
}

.comment-content {
  margin-bottom: 1rem;
}

.comment-text {
  font-size: 1rem;
  line-height: 1.6;
  color: #333;
  word-break: break-word;
}

.comment-text :deep(p) {
  margin: 0.5rem 0;
}

.comment-text :deep(code) {
  background: rgba(102, 126, 234, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
}

.comment-text :deep(pre) {
  background: rgba(102, 126, 234, 0.05);
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 0.5rem 0;
}

.article-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.article-type-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #667eea;
  font-size: 1rem;
}

.article-title {
  font-size: 0.95rem;
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s ease;
  flex: 1;
}

.article-title:hover {
  color: #764ba2;
  text-decoration: underline;
}

.comment-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 0.85rem;
  color: #999;
}

.comment-author {
  font-weight: 500;
  color: #666;
}

.comment-time {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.comment-time svg {
  font-size: 0.8rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-container {
    padding: 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .comment-card {
    padding: 1rem;
  }

  .comment-text {
    font-size: 0.9rem;
  }

  .article-title {
    font-size: 0.85rem;
  }

  .comment-meta {
    font-size: 0.8rem;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}
</style>
