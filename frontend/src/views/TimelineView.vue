<template>
  <div class="timeline-view">
    <NavBar />

    <!-- 日期选择栏 -->
    <div class="date-selector">
      <div class="selector-container">
        <label for="year-select">年份:</label>
        <select id="year-select" v-model="selectedYear" @change="scrollToYear">
          <option value="">全部</option>
          <option v-for="year in availableYears" :key="year" :value="year">{{ year }}</option>
        </select>

        <label for="month-select">月份:</label>
        <select id="month-select" v-model="selectedMonth" @change="scrollToMonth">
          <option value="">全部</option>
          <option v-for="month in availableMonths" :key="month.value" :value="month.value">{{ month.name }}</option>
        </select>
      </div>
    </div>

    <!-- 右侧信息栏 -->
    <div class="right-info-panel">
      <div class="info-content">
        <h3 class="info-title">时间树</h3>
        <p class="info-subtitle">记录生活的点点滴滴</p>
        <div class="info-stats">
          <div class="stat-item">
            <span class="stat-number">{{ totalArticles }}</span>
            <span class="stat-label">文章总数</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ availableYears.length }}</span>
            <span class="stat-label">年份跨度</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 回到顶部按钮 -->
    <BackToTop />

    <div class="timeline-container">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="allArticlesSorted.length === 0" class="empty">
        <div class="empty-icon">📝</div>
        <p>还没有发布任何文章哦～</p>
      </div>

      <div v-else class="timeline-content">
        <!-- 时间树 -->
        <div class="timeline-tree">
          <!-- 主干 -->
          <div class="tree-trunk"></div>

          <!-- 树枝和文章卡片 -->
          <div class="tree-branches">
            <!-- 年份锚点（用于定位，但不显示） -->
            <div
              v-for="year in availableYears"
              :id="`year-${year}`"
              :key="year"
              class="year-anchor"
            ></div>

            <!-- 所有文章按时间排序，左右交替排列 -->
            <div
              v-for="(article, index) in allArticlesSorted"
              :key="`${article.type}-${article.id}`"
              class="article-branch"
              :class="index % 2 === 0 ? 'left-branch' : 'right-branch'"
            >
              <!-- 文章卡片 -->
              <article
                class="blog-card"
                @click="goToArticle(article)"
              >
                <div class="card-content">
                  <div class="card-text">
                    <div class="card-header">
                      <h4 class="card-title">
                        <font-awesome-icon :icon="getTypeIcon(article.type)" class="card-icon" />
                        {{ getArticleTitle(article) }}
                      </h4>
                    </div>
                    <p class="card-excerpt">{{ getArticleExcerpt(article) || '记录成长路上的点点滴滴...' }}</p>
                    <div class="card-meta">
                      <div class="card-date">
                        <i class="date-icon">🕐</i>
                        <span>{{ formatDate(article.createdAt) }}</span>
                      </div>
                      <div class="card-views">
                        <i class="view-icon">👁️</i>
                        <span>{{ article.viewCount || 0 }}</span>
                      </div>
                    </div>
                  </div>
                  <div v-if="article.image" class="card-thumbnail">
                    <img :src="article.image" :alt="getArticleTitle(article)" />
                  </div>
                </div>
              </article>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import BackToTop from '@/components/BackToTop.vue'
import { getArticlesList } from '@/api/Articles/browse'
import { getMomentsList } from '@/api/Moments/browse'

const router = useRouter()

const loading = ref(true)
const timelineData = ref({})
const selectedYear = ref('')
const selectedMonth = ref('')

// 计算属性
const availableYears = computed(() => {
  return Object.keys(timelineData.value).sort((a, b) => b - a)
})

// 所有文章按时间排序（最新的在前）
const allArticlesSorted = computed(() => {
  const articles = []

  // 遍历所有年份和月份，收集所有文章
  Object.values(timelineData.value).forEach(yearData => {
    Object.values(yearData.months).forEach(monthData => {
      articles.push(...monthData.articles)
    })
  })

  // 按创建时间排序（最新的在前）
  return articles.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
})

// 文章总数
const totalArticles = computed(() => {
  return allArticlesSorted.value.length
})

// 获取类型图标
const getTypeIcon = (type) => {
  const iconMap = {
    blog: 'blog',
    project: 'diagram-project',
    research: 'flask',
    moment: 'comment-dots'
  }
  return iconMap[type] || 'file'
}

// 获取文章标题，如果没有标题则使用内容的第一行
const getArticleTitle = (article) => {
  if (article.title && article.title.trim()) {
    return article.title
  }

  // 如果没有标题，从内容中提取第一行
  if (article.content) {
    const firstLine = article.content.split('\n')[0].trim()

    // 移除Markdown语法
    const cleanTitle = firstLine
      .replace(/^#+\s*/, '') // 移除标题标记
      .replace(/\*\*(.*?)\*\*/g, '$1') // 移除粗体标记
      .replace(/\*(.*?)\*/g, '$1') // 移除斜体标记
      .replace(/`(.*?)`/g, '$1') // 移除代码标记
      .replace(/\[(.*?)\]\(.*?\)/g, '$1') // 移除链接标记
      .substring(0, 50) // 限制长度

    return cleanTitle || '无标题'
  }

  return '无标题'
}

// 获取文章摘要
const getArticleExcerpt = (article) => {
  if (article.content) {
    // 移除markdown标记，获取纯文本
    const plainText = article.content
      .replace(/#{1,6}\s+/g, '') // 移除标题标记
      .replace(/\*\*(.*?)\*\*/g, '$1') // 移除粗体标记
      .replace(/\*(.*?)\*/g, '$1') // 移除斜体标记
      .replace(/`(.*?)`/g, '$1') // 移除代码标记
      .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // 移除链接标记
      .replace(/\n+/g, ' ') // 替换换行为空格
      .trim()

    return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText
  }
  return ''
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 可用月份
const availableMonths = computed(() => {
  const months = [
    { value: '01', name: '一月' },
    { value: '02', name: '二月' },
    { value: '03', name: '三月' },
    { value: '04', name: '四月' },
    { value: '05', name: '五月' },
    { value: '06', name: '六月' },
    { value: '07', name: '七月' },
    { value: '08', name: '八月' },
    { value: '09', name: '九月' },
    { value: '10', name: '十月' },
    { value: '11', name: '十一月' },
    { value: '12', name: '十二月' }
  ]
  return months
})

// 滚动到指定年份
const scrollToYear = () => {
  if (selectedYear.value) {
    const element = document.getElementById(`year-${selectedYear.value}`)
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  }
}

// 滚动到指定月份
const scrollToMonth = () => {
  if (selectedYear.value && selectedMonth.value) {
    // 这里可以根据需要实现更精确的月份定位
    scrollToYear()
  }
}

// 跳转到文章详情
const goToArticle = (article) => {
  router.push(`/${article.type}/${article.id}`)
}

// 加载所有文章并构建时间线
const loadTimelineData = async () => {
  try {
    loading.value = true
    const types = ['blog', 'project', 'research', 'moment']
    const allArticles = []

    // 加载所有类型的文章
    for (const type of types) {
      if (type === 'moment') {
        const response = await getMomentsList(1, 1000)
        const articles = response.data.map(item => ({
          id: item.ID,
          title: item.Title, // 注意：后端返回的是 Title（大写）
          content: item.Content, // 注意：后端返回的是 Content（大写）
          type: 'moment',
          image: item.Image, // 注意：后端返回的是 Image（大写）
          tags: [],
          viewCount: item.viewCount || 0, // 注意：后端返回的是 viewCount
          createdAt: item.CreatedAt // 注意：后端返回的是 CreatedAt（大写）
        }))
        allArticles.push(...articles)
      } else {
        const response = await getArticlesList(type, 1, 1000)
        const articles = response.data.map(item => ({
          id: item.ID,
          title: item.title,
          content: item.content || item.abstract || item.status || '', // 根据类型使用不同字段
          type,
          image: item.image,
          tags: item.tags || [],
          viewCount: item.viewCount || 0, // 注意：GORM返回的是 viewCount
          createdAt: item.CreatedAt // 注意：GORM返回的是 CreatedAt（大写）
        }))
        allArticles.push(...articles)
      }
    }

    // 按时间排序（最新的在前）
    allArticles.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))

    // 构建时间线数据结构
    const timeline = {}
    allArticles.forEach(article => {
      const date = new Date(article.createdAt)
      const year = date.getFullYear().toString()
      const month = (date.getMonth() + 1).toString().padStart(2, '0')

      if (!timeline[year]) {
        timeline[year] = {
          total: 0,
          months: {}
        }
      }

      if (!timeline[year].months[month]) {
        timeline[year].months[month] = {
          articles: []
        }
      }

      timeline[year].months[month].articles.push(article)
      timeline[year].total++
    })
    timelineData.value = timeline
  } catch (error) {
    console.error('加载时间线数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTimelineData()
})
</script>

<style scoped>
.timeline-view {
  min-height: 100vh;
  padding-top: 80px;
  position: relative;
  overflow-x: hidden;
}

.timeline-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px 60px;
  position: relative;
  z-index: 1;
}

.loading,
.empty {
  text-align: center;
  padding: 4rem 2rem;
  color: #666;
  font-size: 1.1rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  margin: 2rem 0;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top: 4px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.empty-icon {
  font-size: 4rem;
  opacity: 0.7;
}

/* 年份导航 */
.year-nav {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 3rem;
  padding: 1.5rem;
  background: transparent;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.year-btn {
  padding: 0.8rem 1.5rem;
  background: #f5f5f5;
  border: none;
  border-radius: 25px;
  color: #666;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.year-btn:hover {
  background: #e0e0e0;
  color: #333;
}

.year-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* 日期选择器 */
.date-selector {
  position: fixed;
  top: 200px;
  left: 20px;
  z-index: 100;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(15px);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.date-selector:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}

/* 右侧信息栏 */
.right-info-panel {
  position: fixed;
  top: 200px;
  right: 20px;
  z-index: 100;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(15px);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
  max-width: 250px;
}

.right-info-panel:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}

.info-content {
  text-align: center;
}

.info-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.info-subtitle {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 20px;
  line-height: 1.4;
}

.info-stats {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.stat-number {
  font-size: 1.8rem;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.8rem;
  color: #666;
  font-weight: 500;
}

.selector-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.selector-container label {
  font-size: 0.9rem;
  color: #6a1b9a;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.selector-container select {
  padding: 10px 15px;
  border: 2px solid rgba(106, 27, 154, 0.2);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
  color: #333;
  transition: all 0.3s ease;
  cursor: pointer;
}

.selector-container select:focus {
  outline: none;
  border-color: #6a1b9a;
  box-shadow: 0 0 0 3px rgba(106, 27, 154, 0.1);
}

/* 时间树样式 */
.timeline-tree {
  position: relative;
  max-width: 1400px;
  margin: 0 auto;
  padding: 60px 0;
  margin-top: 40px;
}

.timeline-tree::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(ellipse at center,
    rgba(102, 126, 234, 0.05) 0%,
    rgba(118, 75, 162, 0.03) 30%,
    transparent 70%
  );
  transform: translateX(-50%);
  z-index: 0;
  pointer-events: none;
}

/* 主干 */
.tree-trunk {
  position: absolute;
  left: 50%;
  top: 60px;
  bottom: 0;
  width: 8px;
  background: linear-gradient(to bottom,
    rgba(255, 255, 255, 0.9) 0%,
    rgba(102, 126, 234, 0.9) 15%,
    rgba(118, 75, 162, 0.9) 35%,
    rgba(102, 126, 234, 0.9) 55%,
    rgba(118, 75, 162, 0.9) 75%,
    rgba(102, 126, 234, 0.9) 90%,
    rgba(255, 255, 255, 0.9) 100%
  );
  transform: translateX(-50%);
  z-index: 1;
  border-radius: 4px;
  box-shadow:
    0 0 30px rgba(102, 126, 234, 0.4),
    0 0 60px rgba(118, 75, 162, 0.2),
    inset 0 0 10px rgba(255, 255, 255, 0.3);
  animation: trunkGlow 3s ease-in-out infinite alternate;
}

@keyframes trunkGlow {
  0% {
    box-shadow:
      0 0 30px rgba(102, 126, 234, 0.4),
      0 0 60px rgba(118, 75, 162, 0.2),
      inset 0 0 10px rgba(255, 255, 255, 0.3);
  }
  100% {
    box-shadow:
      0 0 40px rgba(102, 126, 234, 0.6),
      0 0 80px rgba(118, 75, 162, 0.4),
      inset 0 0 15px rgba(255, 255, 255, 0.5);
  }
}

/* 方案1: 简洁的圆点装饰 - 当前使用 */
.tree-trunk::before { content: none; }

@keyframes topGlow {
  0% {
    box-shadow: 0 0 20px rgba(102, 126, 234, 0.4);
    transform: scale(1);
  }
  100% {
    box-shadow: 0 0 30px rgba(102, 126, 234, 0.6);
    transform: scale(1.1);
  }
}

.tree-trunk::after { content: none; }

@keyframes bottomGlow {
  0% {
    box-shadow: 0 0 20px rgba(118, 75, 162, 0.4);
    transform: scale(1);
  }
  100% {
    box-shadow: 0 0 30px rgba(118, 75, 162, 0.6);
    transform: scale(1.1);
  }
}

/* 树枝区域 */
.tree-branches {
  position: relative;
  z-index: 2;
}

/* 年份区域 */
.tree-year-section {
  margin-bottom: 80px;
  position: relative;
}

/* 年份锚点（用于定位但不显示） */
.year-anchor {
  position: absolute;
  left: 50%;
  top: 0;
  width: 1px;
  height: 1px;
  transform: translateX(-50%);
  z-index: 1;
}

/* 文章分支 */
.article-branch {
  position: relative;
  margin: 40px 0;
  padding: 0 40px;
}

/* 左分支 */
.left-branch {
  text-align: left;
}

.left-branch::before {
  content: '';
  position: absolute;
  left: calc(50% - 120px);
  top: 50%;
  width: 120px;
  height: 3px;
  background: linear-gradient(to left,
    rgba(102, 126, 234, 0.8) 0%,
    rgba(118, 75, 162, 0.6) 30%,
    rgba(102, 126, 234, 0.4) 60%,
    transparent 100%
  );
  z-index: 1;
  transform: translateY(-50%);
  border-radius: 2px;
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.3);
  animation: branchGlow 4s ease-in-out infinite alternate;
}

@keyframes branchGlow {
  0% {
    box-shadow: 0 0 10px rgba(102, 126, 234, 0.3);
  }
  100% {
    box-shadow: 0 0 20px rgba(102, 126, 234, 0.5);
  }
}

.left-branch::after { content: none; }

@keyframes connectionPulse {
  0% {
    transform: translateY(-50%) scale(1);
    box-shadow:
      0 0 15px rgba(102, 126, 234, 0.6),
      0 0 30px rgba(118, 75, 162, 0.3),
      inset 0 0 4px rgba(255, 255, 255, 0.8);
  }
  100% {
    transform: translateY(-50%) scale(1.1);
    box-shadow:
      0 0 20px rgba(102, 126, 234, 0.8),
      0 0 40px rgba(118, 75, 162, 0.5),
      inset 0 0 6px rgba(255, 255, 255, 1);
  }
}

/* 右分支 */
.right-branch {
  text-align: right;
}

.right-branch::before {
  content: '';
  position: absolute;
  right: calc(50% - 120px);
  top: 50%;
  width: 120px;
  height: 3px;
  background: linear-gradient(to right,
    rgba(102, 126, 234, 0.8) 0%,
    rgba(118, 75, 162, 0.6) 30%,
    rgba(102, 126, 234, 0.4) 60%,
    transparent 100%
  );
  z-index: 1;
  transform: translateY(-50%);
  border-radius: 2px;
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.3);
  animation: branchGlow 4s ease-in-out infinite alternate;
}

.right-branch::after {
  content: none;
}

/* 月份标签 */
.month-label {
  font-size: 1rem;
  color: #666;
  margin-bottom: 15px;
  font-weight: 500;
}

/* 文章容器 */
.articles-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* 博客卡片样式 - 与主页保持一致 */
.blog-card {
  background: transparent;
  border-radius: 15px;
  border: none;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  max-width: 400px;
  margin: 0 auto;
}

/* 左右分支布局 */
.left-branch .blog-card {
  margin-right: auto;
  margin-left: 0;
  text-align: left;
}

.right-branch .blog-card {
  margin-left: auto;
  margin-right: 0;
  text-align: right;
}

.blog-card:hover {
  background: linear-gradient(135deg, rgba(248, 249, 255, 0.3) 0%, rgba(240, 242, 255, 0.3) 100%);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.15);
}

.card-content {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 20px;
  min-height: 100px;
  border-radius: 15px;
  transition: all 0.3s ease;
}

/* 左分支：文字在左，图片在右 */
.left-branch .card-content {
  flex-direction: row;
}

/* 右分支：图片在左，文字在右 */
.right-branch .card-content {
  flex-direction: row-reverse;
}

.blog-card:hover .card-content {
  background: rgba(255, 255, 255, 0.05);
}

.card-text {
  flex: 0 0 65%;
  display: flex;
  flex-direction: column;
  gap: 12px;
  text-align: left;
  justify-content: space-between;
}

.right-branch .card-text {
  text-align: right;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.right-branch .card-header {
  justify-content: flex-end;
}

.card-icon {
  font-size: 1rem;
  margin-right: 8px;
  color: #667eea;
  vertical-align: middle;
}

.card-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: #333;
  line-height: 1.4;
  margin: 0;
  text-align: left;
}

.card-excerpt {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.right-branch .card-meta {
  flex-direction: row-reverse;
}

.card-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.85rem;
  color: #666;
  font-weight: 500;
}

.date-icon {
  font-size: 0.8rem;
}

.card-views {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.85rem;
  color: #666;
  font-weight: 500;
}

.view-icon {
  font-size: 0.8rem;
}

/* 图片缩略图样式 */
.card-thumbnail {
  flex: 0 0 35%;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
}

.card-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.blog-card:hover .card-thumbnail img {
  transform: scale(1.05);
}

@keyframes cardFloat {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-2px);
  }
}

.article-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(106, 27, 154, 0.05) 0%, rgba(142, 45, 226, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 1;
}

.article-card:hover::before {
  opacity: 1;
}

.left-branch .article-card {
  margin-right: auto;
  margin-left: 0;
  text-align: left;
}

.right-branch .article-card {
  margin-left: auto;
  margin-right: 0;
  text-align: right;
}

.article-card:hover {
  transform: translateY(-12px) scale(1.03);
  box-shadow:
    0 20px 60px rgba(0, 0, 0, 0.15),
    0 8px 30px rgba(102, 126, 234, 0.2),
    0 4px 15px rgba(118, 75, 162, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
  background: linear-gradient(135deg,
    rgba(255, 255, 255, 1) 0%,
    rgba(248, 250, 252, 1) 50%,
    rgba(240, 242, 255, 1) 100%
  );
  animation: none;
}

.card-hover-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 2;
  pointer-events: none;
}

.article-card:hover .card-hover-effect {
  opacity: 1;
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 12px;
  position: relative;
  z-index: 3;
}

.card-type {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border: 1px solid rgba(102, 126, 234, 0.2);
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.card-type:hover {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
  transform: scale(1.05);
}

.type-icon {
  font-size: 1.2rem;
  color: #667eea;
  transition: all 0.3s ease;
}

.card-type:hover .type-icon {
  color: #5a67d8;
  transform: scale(1.1);
}

.card-date {
  font-size: 0.8rem;
  color: #666;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.8);
  padding: 4px 8px;
  border-radius: 8px;
  backdrop-filter: blur(10px);
}

.card-views {
  font-size: 0.75rem;
  color: #6a1b9a;
  font-weight: 500;
  background: rgba(106, 27, 154, 0.1);
  padding: 3px 8px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 4px;
  border: 1px solid rgba(106, 27, 154, 0.15);
  white-space: nowrap;
  flex-shrink: 0;
}

.view-icon {
  font-size: 0.7rem;
}

/* 卡片标题 */
.card-title {
  font-size: 1.3rem;
  font-weight: 700;
  background: linear-gradient(135deg,
    #2c3e50 0%,
    #34495e 30%,
    #667eea 60%,
    #764ba2 100%
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  position: relative;
  z-index: 3;
  letter-spacing: 0.3px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.article-card:hover .card-title {
  background: linear-gradient(135deg,
    #667eea 0%,
    #764ba2 50%,
    #667eea 100%
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 卡片标签 */

/* 卡片底部 */
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(106, 27, 154, 0.1);
  position: relative;
  z-index: 3;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .timeline-container {
    padding: 0 40px 60px;
  }

  .article-card {
    max-width: 320px;
  }
}

@media (max-width: 768px) {

  .date-selector {
    position: fixed;
    top: 180px;
    left: 20px;
    z-index: 100;
    transform: none;
    margin: 0;
  }

  .right-info-panel {
    position: fixed;
    top: 180px;
    right: 20px;
    z-index: 100;
    max-width: 200px;
    padding: 15px;
  }

  .selector-container {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 15px;
  }

  .timeline-tree {
    padding: 40px 0;
  }

  .tree-trunk {
    left: 30px;
    width: 4px;
  }

  .article-branch {
    padding: 0 60px 0 20px;
  }

  .left-branch::before {
    left: -30px;
    width: 30px;
  }

  .left-branch::after {
    left: -30px;
  }

  .right-branch::before {
    right: -30px;
    width: 30px;
  }

  .right-branch::after {
    right: -30px;
  }

  .left-branch .article-card,
  .right-branch .article-card {
    margin: 0;
    max-width: 100%;
    padding: 15px;
  }

  .card-title {
    font-size: 1.1rem;
  }

  .card-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}

@media (max-width: 480px) {

  .timeline-container {
    padding: 0 15px 40px;
  }

  .date-selector {
    position: fixed;
    top: 160px;
    left: 10px;
    z-index: 100;
    padding: 15px;
    max-width: 200px;
  }

  .right-info-panel {
    position: fixed;
    top: 160px;
    right: 10px;
    z-index: 100;
    max-width: 180px;
    padding: 12px;
  }

  .info-title {
    font-size: 1.2rem;
  }

  .stat-number {
    font-size: 1.5rem;
  }

  .article-branch {
    padding: 0 40px 0 15px;
  }

  .tree-trunk {
    left: 20px;
  }

  .left-branch::before,
  .right-branch::before {
    width: 20px;
  }

  .left-branch::before {
    left: -20px;
  }

  .left-branch::after {
    left: -20px;
  }

  .right-branch::before {
    right: -20px;
  }

  .right-branch::after {
    right: -20px;
  }
}

/* 时间线 */
.timeline {
  background: transparent;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.timeline-year {
  margin-bottom: 3rem;
}

.timeline-year:last-child {
  margin-bottom: 0;
}

/* ========== 其他时间树设计方案 ========== */

/* 方案2: 完全移除顶部和底部装饰 - 最简洁 */
/*
.tree-trunk::before,
.tree-trunk::after {
  display: none;
}
*/

/* 方案3: 叶子形状装饰 */
/*
.tree-trunk::before {
  content: '';
  position: absolute;
  top: -20px;
  left: -10px;
  width: 0;
  height: 0;
  border-left: 15px solid transparent;
  border-right: 15px solid transparent;
  border-bottom: 25px solid rgba(102, 126, 234, 0.8);
  filter: drop-shadow(0 0 10px rgba(102, 126, 234, 0.4));
  animation: leafSway 4s ease-in-out infinite alternate;
}

.tree-trunk::after {
  content: '';
  position: absolute;
  bottom: -20px;
  left: -10px;
  width: 0;
  height: 0;
  border-left: 15px solid transparent;
  border-right: 15px solid transparent;
  border-top: 25px solid rgba(118, 75, 162, 0.8);
  filter: drop-shadow(0 0 10px rgba(118, 75, 162, 0.4));
  animation: leafSway 4s ease-in-out infinite alternate;
}

@keyframes leafSway {
  0% { transform: rotate(-5deg); }
  100% { transform: rotate(5deg); }
}
*/

/* 方案4: 星星装饰 */
/*
.tree-trunk::before {
  content: '⭐';
  position: absolute;
  top: -20px;
  left: -12px;
  font-size: 24px;
  animation: starTwinkle 2s ease-in-out infinite alternate;
}

.tree-trunk::after {
  content: '✨';
  position: absolute;
  bottom: -20px;
  left: -12px;
  font-size: 24px;
  animation: starTwinkle 2s ease-in-out infinite alternate;
}

@keyframes starTwinkle {
  0% {
    transform: scale(1) rotate(0deg);
    filter: brightness(1);
  }
  100% {
    transform: scale(1.2) rotate(180deg);
    filter: brightness(1.5);
  }
}
*/

/* 方案5: 渐变线条装饰 */
/*
.tree-trunk::before {
  content: '';
  position: absolute;
  top: -25px;
  left: -2px;
  width: 4px;
  height: 25px;
  background: linear-gradient(to bottom,
    rgba(102, 126, 234, 0.8) 0%,
    rgba(118, 75, 162, 0.6) 50%,
    transparent 100%
  );
  border-radius: 2px;
  box-shadow: 0 0 15px rgba(102, 126, 234, 0.3);
}

.tree-trunk::after {
  content: '';
  position: absolute;
  bottom: -25px;
  left: -2px;
  width: 4px;
  height: 25px;
  background: linear-gradient(to top,
    rgba(118, 75, 162, 0.8) 0%,
    rgba(102, 126, 234, 0.6) 50%,
    transparent 100%
  );
  border-radius: 2px;
  box-shadow: 0 0 15px rgba(118, 75, 162, 0.3);
}
*/

.year-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 3px solid #f0f0f0;
}

.year-title {
  font-size: 2rem;
  color: #333;
  margin: 0;
  font-weight: 700;
}

.year-count {
  color: #999;
  font-size: 1rem;
  background: #f5f5f5;
  padding: 0.5rem 1rem;
  border-radius: 15px;
}

.months-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.timeline-month {
  border-left: 3px solid #e0e0e0;
  padding-left: 2rem;
  position: relative;
}

.timeline-month::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 0;
  width: 12px;
  height: 12px;
  background: #667eea;
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 0 0 3px #e0e0e0;
}

.month-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.month-title {
  font-size: 1.3rem;
  color: #333;
  margin: 0;
  font-weight: 600;
}

.month-count {
  color: #999;
  font-size: 0.9rem;
  background: #f5f5f5;
  padding: 0.3rem 0.8rem;
  border-radius: 10px;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.article-card {
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  background: #f9f9f9;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.article-card:hover {
  background: #f0f0f0;
  border-color: #667eea;
  transform: translateY(-2px);
}

.article-date {
  flex-shrink: 0;
  width: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.article-date .day {
  font-size: 1.5rem;
  font-weight: 700;
  color: #667eea;
  background: transparent;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.article-content {
  flex: 1;
}

.article-type {
  display: inline-block;
  padding: 0.3rem 0.8rem;
  background: #667eea;
  color: white;
  border-radius: 12px;
  font-size: 0.8rem;
  margin-bottom: 0.8rem;
}

.article-title {
  font-size: 1.2rem;
  color: #333;
  margin: 0 0 0.8rem 0;
  font-weight: 600;
  line-height: 1.4;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.article-tags .tag {
  padding: 0.2rem 0.6rem;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 10px;
  font-size: 0.8rem;
}

.more-tags {
  padding: 0.2rem 0.6rem;
  background: #f5f5f5;
  color: #999;
  border-radius: 10px;
  font-size: 0.8rem;
}

.article-image {
  flex-shrink: 0;
  width: 100px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

@media (max-width: 768px) {
  .timeline-container {
    padding: 1rem;
  }

  .timeline-stats {
    flex-direction: column;
    gap: 0.5rem;
  }

  .year-nav {
    padding: 1rem;
  }

  .timeline {
    padding: 1rem;
  }

  .timeline-month {
    padding-left: 1rem;
  }

  .article-card {
    flex-direction: column;
    gap: 1rem;
  }

  .article-date {
    width: auto;
    justify-content: flex-start;
  }

  .article-image {
    width: 100%;
    height: 150px;
  }
}
</style>
