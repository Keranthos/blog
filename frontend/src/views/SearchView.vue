<template>
  <div class="search-view">
    <NavBar />

    <!-- 搜索头部 -->
    <div class="search-header">
      <div class="search-container">
        <h1 class="search-title">搜索结果</h1>
        <div class="search-input-wrapper">
          <div class="search-input-container">
            <input
              ref="searchInput"
              v-model="searchKeyword"
              type="text"
              placeholder="搜索文章..."
              class="search-input"
              @keydown.enter="performSearch"
              @input="onInputChange"
              @focus="showSuggestions = true"
              @blur="onInputBlur"
            />
            <SearchSuggestions
              :show="showSuggestions && !searchKeyword"
              :keyword="searchKeyword"
              :search-history="searchHistory"
              @select="onSuggestionSelect"
              @remove-from-history="removeFromHistory"
              @close="showSuggestions = false"
            />
          </div>
          <button class="search-btn" @click="performSearch">
            <font-awesome-icon icon="magnifying-glass" />
            搜索
          </button>
        </div>
        <div v-if="searchKeyword" class="search-info">
          搜索关键词：<span class="keyword">{{ searchKeyword }}</span>
          <span v-if="!loading && searchResults.length > 0" class="result-count">
            找到 {{ searchResults.length }} 个结果
          </span>
        </div>

        <!-- 搜索历史 -->
        <div v-if="!searchKeyword && searchHistory.length > 0" class="search-history">
          <div class="history-header">
            <h3>搜索历史</h3>
            <button class="clear-btn" @click="clearHistory">清空</button>
          </div>
          <div class="history-tags">
            <span
              v-for="(item, index) in searchHistory"
              :key="index"
              class="history-tag"
              @click="searchFromHistory(item)"
            >
              {{ item }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索结果 -->
    <div class="search-results-container">
      <!-- 加载中 -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>搜索中...</p>
      </div>

      <!-- 搜索结果列表 -->
      <div v-else-if="searchResults.length > 0" class="results-list">
        <div
          v-for="result in searchResults"
          :key="result.id"
          class="result-item"
          @click="goToArticle(result)"
        >
          <div class="result-image">
            <img :src="result.image || 'https://picsum.photos/id/1/400/300'" alt="封面" />
          </div>
          <div class="result-content">
            <div class="result-type">{{ getTypeName(result.type) }}</div>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <h3 class="result-title" v-html="highlightKeyword(result.title)"></h3>
            <div v-if="result.tags && result.tags.length" class="result-tags">
              <span v-for="tag in result.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
            <!-- eslint-disable-next-line vue/no-v-html -->
            <p v-if="result.excerpt" class="result-excerpt" v-html="highlightKeyword(result.excerpt)"></p>
            <div class="result-meta">
              <span class="meta-item">
                <font-awesome-icon icon="calendar" />
                {{ formatDate(result.time) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 无结果 -->
      <div v-else-if="!loading && searchKeyword" class="empty-state">
        <font-awesome-icon icon="magnifying-glass" class="empty-icon" />
        <p>没有找到相关内容</p>
        <p class="empty-tip">试试其他关键词吧</p>
      </div>

      <!-- 初始状态 -->
      <div v-else class="empty-state">
        <font-awesome-icon icon="magnifying-glass" class="empty-icon" />
        <p>输入关键词开始搜索</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DOMPurify from 'dompurify'
import NavBar from '@/components/NavBar.vue'
import SearchSuggestions from '@/components/SearchSuggestions.vue'
import { getArticlesList } from '@/api/Articles/browse'
import { getMomentsList } from '@/api/Moments/browse'

const route = useRoute()
const router = useRouter()

const searchKeyword = ref('')
const searchResults = ref([])
const loading = ref(false)
const searchHistory = ref([])
const showSuggestions = ref(false)
const searchInput = ref(null)

// 搜索历史相关
const SEARCH_HISTORY_KEY = 'search_history'
const MAX_HISTORY_LENGTH = 10

// 加载搜索历史
const loadSearchHistory = () => {
  try {
    const history = localStorage.getItem(SEARCH_HISTORY_KEY)
    if (history) {
      searchHistory.value = JSON.parse(history)
    }
  } catch (error) {
    searchHistory.value = []
  }
}

// 保存搜索历史
const saveSearchHistory = (keyword) => {
  if (!keyword || !keyword.trim()) return

  // 移除重复项
  searchHistory.value = searchHistory.value.filter(item => item !== keyword)

  // 添加到开头
  searchHistory.value.unshift(keyword)

  // 限制数量
  if (searchHistory.value.length > MAX_HISTORY_LENGTH) {
    searchHistory.value = searchHistory.value.slice(0, MAX_HISTORY_LENGTH)
  }

  // 保存到 localStorage
  localStorage.setItem(SEARCH_HISTORY_KEY, JSON.stringify(searchHistory.value))
}

// 清空搜索历史
const clearHistory = () => {
  searchHistory.value = []
  localStorage.removeItem(SEARCH_HISTORY_KEY)
}

// 从历史记录搜索
const searchFromHistory = (keyword) => {
  searchKeyword.value = keyword
  performSearch()
}

// 搜索建议相关方法
const onInputChange = () => {
  // 输入时显示建议
  if (!searchKeyword.value.trim()) {
    showSuggestions.value = true
  } else {
    showSuggestions.value = false
  }
}

const onInputBlur = () => {
  // 延迟隐藏建议，以便点击建议项
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}

const onSuggestionSelect = (suggestion) => {
  searchKeyword.value = suggestion
  showSuggestions.value = false
  nextTick(() => {
    performSearch()
  })
}

// 获取类型名称
const getTypeName = (type) => {
  const typeMap = {
    blog: '博客',
    project: '项目',
    research: '科研',
    moment: '碎碎念'
  }
  return typeMap[type] || type
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

// 高亮关键词（使用 DOMPurify 防止 XSS）
const highlightKeyword = (text) => {
  if (!text || !searchKeyword.value) return text
  const keyword = searchKeyword.value.trim()
  // 转义特殊字符，防止正则注入
  const escapedKeyword = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(`(${escapedKeyword})`, 'gi')
  const highlighted = text.replace(regex, '<mark>$1</mark>')
  // 使用 DOMPurify 清理，只允许 <mark> 标签
  return DOMPurify.sanitize(highlighted, {
    ALLOWED_TAGS: ['mark'],
    ALLOWED_ATTR: []
  })
}

// 执行搜索
const performSearch = async () => {
  const keyword = searchKeyword.value.trim()
  if (!keyword) return

  // 保存到搜索历史
  saveSearchHistory(keyword)

  loading.value = true
  searchResults.value = []

  try {
    // 搜索所有类型的文章
    const types = ['blog', 'project', 'research', 'moment']
    const allResults = []

    for (const type of types) {
      try {
        let articles = []
        if (type === 'moment') {
          const response = await getMomentsList(1, 100) // 获取更多结果用于搜索
          articles = response.data.map(item => ({
            id: item.ID,
            type: 'moment',
            title: item.Title,
            content: item.Content,
            image: item.Image,
            tags: [],
            time: item.CreatedAt,
            excerpt: item.Content ? item.Content.substring(0, 150) + '...' : ''
          }))
        } else {
          const response = await getArticlesList(type, 1, 100)
          articles = response.data.map(item => ({
            id: item.ID,
            type,
            title: item.title,
            content: item.content || item.abstract || '',
            image: item.image,
            tags: item.tags || [],
            time: item.CreatedAt,
            excerpt: (item.content || item.abstract || '').substring(0, 150) + '...'
          }))
        }

        // 过滤匹配关键词的文章
        const filtered = articles.filter(article => {
          const titleMatch = article.title.toLowerCase().includes(keyword.toLowerCase())
          const contentMatch = article.content.toLowerCase().includes(keyword.toLowerCase())
          const tagsMatch = article.tags.some(tag => tag.toLowerCase().includes(keyword.toLowerCase()))
          return titleMatch || contentMatch || tagsMatch
        })

        allResults.push(...filtered)
      } catch (error) {
        console.error(`搜索 ${type} 失败:`, error)
      }
    }

    searchResults.value = allResults
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    loading.value = false
  }
}

// 跳转到文章详情
const goToArticle = (result) => {
  if (result.type === 'moment') {
    router.push(`/moments/${result.id}`)
  } else if (result.type === 'blog') {
    router.push({ name: 'BlogDetail', params: { id: result.id } })
  } else if (result.type === 'research' || result.type === 'project') {
    router.push({ name: 'BlogDetail', params: { id: result.id }, query: { type: result.type } })
  }
}

// 初始化
onMounted(() => {
  loadSearchHistory()

  const keyword = route.query.search
  if (keyword) {
    searchKeyword.value = keyword
    performSearch()
  }
})
</script>

<style scoped>
.search-view {
  min-height: 100vh;
  padding-top: 100px;
  padding-bottom: 80px;
}

/* 搜索头部 */
.search-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60px 20px 80px;
  margin-bottom: -40px;
}

.search-container {
  max-width: 800px;
  margin: 0 auto;
}

.search-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  text-align: center;
  margin-bottom: 30px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.search-input-wrapper {
  display: flex;
  gap: 10px;
  background: transparent;
  border-radius: 50px;
  padding: 8px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.search-input-container {
  position: relative;
  flex: 1;
}

.search-input {
  flex: 1;
  padding: 12px 20px;
  border: none;
  background: transparent;
  font-size: 1rem;
  outline: none;
}

.search-btn {
  padding: 12px 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 50px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.search-info {
  margin-top: 20px;
  text-align: center;
  color: white;
  font-size: 1rem;
}

.keyword {
  font-weight: 700;
  color: #ffd700;
  padding: 0 8px;
}

.result-count {
  margin-left: 15px;
  opacity: 0.9;
}

/* 搜索历史 */
.search-history {
  margin-top: 2rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.history-header h3 {
  font-size: 1.1rem;
  color: #333;
  margin: 0;
}

.clear-btn {
  padding: 0.4rem 1rem;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  color: #666;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.clear-btn:hover {
  background: #e0e0e0;
  color: #333;
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem;
}

.history-tag {
  padding: 0.5rem 1rem;
  background: #f0f0f0;
  border-radius: 20px;
  color: #666;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.history-tag:hover {
  background: #667eea;
  color: white;
  transform: translateY(-2px);
}

/* 搜索结果容器 */
.search-results-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 80px 20px;
  color: #666;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 结果列表 */
.results-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.result-item {
  display: flex;
  gap: 20px;
  background: transparent;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.result-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  border-color: rgba(102, 126, 234, 0.3);
}

.result-image {
  flex-shrink: 0;
  width: 200px;
  height: 150px;
  border-radius: 12px;
  overflow: hidden;
}

.result-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.result-item:hover .result-image img {
  transform: scale(1.1);
}

.result-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.result-type {
  display: inline-block;
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 600;
  width: fit-content;
}

.result-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin: 0;
  line-height: 1.4;
}

.result-title :deep(mark) {
  background: #ffd700;
  color: #333;
  padding: 2px 4px;
  border-radius: 3px;
  font-weight: 700;
}

.result-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tag {
  padding: 4px 12px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.result-excerpt {
  color: #666;
  line-height: 1.6;
  font-size: 0.95rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.result-excerpt :deep(mark) {
  background: #ffd700;
  color: #333;
  padding: 2px 4px;
  border-radius: 3px;
  font-weight: 600;
}

.result-meta {
  display: flex;
  gap: 15px;
  color: #999;
  font-size: 0.9rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 100px 20px;
  color: #999;
}

.empty-icon {
  font-size: 5rem;
  color: #ddd;
  margin-bottom: 20px;
}

.empty-state p {
  font-size: 1.2rem;
  margin-bottom: 10px;
}

.empty-tip {
  font-size: 1rem;
  color: #bbb;
}

/* 响应式 */
@media (max-width: 768px) {
  .search-view {
    padding-top: 80px;
  }

  .search-title {
    font-size: 2rem;
  }

  .result-item {
    flex-direction: column;
  }

  .result-image {
    width: 100%;
    height: 200px;
  }
}
</style>
