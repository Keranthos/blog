<template>
  <ModernLoading
    v-if="!isComponentReady"
    :progress="loadingProgress"
    :title="'山角函兽的小窝'"
    :subtitle="'Loading……'"
  />
  <div v-else class="home-view">
    <NavBar />
    <div class="hero-section">
      <div class="hero-content">
        <div class="avatar-container">
          <img src="@/assets/my_headportrait.jpg" alt="Avatar" class="hero-avatar" />
          <div class="avatar-ring"></div>
        </div>
        <h1 class="hero-title">{{ displayedText }}<span class="cursor">|</span></h1>
        <p class="hero-subtitle">记录生活 · 分享技术 · 探索未知</p>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">{{ blogCount }}</span>
            <span class="stat-label">博客文章</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ momentCount }}</span>
            <span class="stat-label">随笔</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">{{ projectCount }}</span>
            <span class="stat-label">项目</span>
          </div>
        </div>
      </div>
      <div class="scroll-indicator">
        <span>向下滚动</span>
        <div class="scroll-arrow">↓</div>
      </div>
    </div>
    <!-- 现代化内容区域 - 模仿kirigaya.cn布局 -->
    <div class="modern-content">
      <div class="main-layout">
        <!-- 主内容区域 -->
        <div class="main-content">
          <!-- 最新博客 -->
          <div class="content-section">
            <div class="section-header">
              <h3 class="section-title">
                <span class="section-icon">📝</span>
                <span class="section-text">博客</span>
              </h3>
              <div class="section-count">{{ blogCount }} 篇文章</div>
            </div>
            <div class="blog-cards">
              <article
                v-for="blog in latestBlogs"
                :key="blog.ID"
                class="blog-card"
                @click="$router.push(`/blog/${blog.ID}`)"
              >
                <div class="card-content">
                  <div class="card-text">
                    <div class="card-header">
                      <h4 class="card-title">
                        <font-awesome-icon icon="blog" class="card-icon" />
                        {{ blog.title }}
                      </h4>
                    </div>
                    <p class="card-excerpt">{{ getPlainText(blog.content) || '记录成长路上的点点滴滴...' }}</p>
                  </div>
                  <div class="card-thumbnail">
                    <img :src="blog.image" :alt="blog.title" />
                  </div>
                  <div class="card-meta">
                    <div class="card-date">
                      <i class="date-icon">🕐</i>
                      <span>{{ formatDateTime(blog.CreatedAt) }}</span>
                    </div>
                    <div class="card-tags">
                      <i class="tag-icon">🔖</i>
                      <template v-if="blog.tags && blog.tags.length > 0">
                        <span v-for="tag in blog.tags" :key="tag" class="tag">{{ tag }}</span>
                      </template>
                      <span v-else class="tag">山角函兽懒得加标签了</span>
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </div>

          <!-- 最新随笔 -->
          <div class="content-section">
            <div class="section-header">
              <h3 class="section-title">
                <span class="section-icon">✍️</span>
                <span class="section-text">随笔</span>
              </h3>
              <div class="section-count">{{ momentCount }} 篇随笔</div>
            </div>
            <div class="blog-cards">
              <article
                v-for="moment in latestMoments"
                :key="moment.ID"
                class="blog-card"
                @click="$router.push(`/moments/${moment.ID}`)"
              >
                <div class="card-content">
                  <div class="card-text">
                    <div class="card-header">
                      <h4 class="card-title">
                        <font-awesome-icon icon="pen-to-square" class="card-icon" />
                        {{ moment.title }}
                      </h4>
                    </div>
                    <p class="card-excerpt">{{ getPlainText(moment.content) || '记录生活的点点滴滴...' }}</p>
                  </div>
                  <div class="card-thumbnail">
                    <img :src="moment.image" :alt="moment.title" />
                  </div>
                  <div class="card-meta">
                    <div class="card-date">
                      <i class="date-icon">🕐</i>
                      <span>{{ formatDateTime(moment.CreatedAt) }}</span>
                    </div>
                    <div class="card-tags">
                      <i class="tag-icon">🔖</i>
                      <template v-if="moment.tags && moment.tags.length > 0">
                        <span v-for="tag in moment.tags" :key="tag" class="tag">{{ tag }}</span>
                      </template>
                      <span v-else class="tag">山角函兽懒得加标签了</span>
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </div>

          <!-- 最新项目 -->
          <div class="content-section">
            <div class="section-header">
              <h3 class="section-title">
                <span class="section-icon">💼</span>
                <span class="section-text">项目</span>
              </h3>
              <div class="section-count">{{ projectCount }} 个项目</div>
            </div>
            <div class="blog-cards">
              <article
                v-for="project in latestProjects"
                :key="project.ID"
                class="blog-card"
                @click="$router.push(`/blog/${project.ID}`)"
              >
                <div class="card-content">
                  <div class="card-text">
                    <div class="card-header">
                      <h4 class="card-title">
                        <font-awesome-icon icon="diagram-project" class="card-icon" />
                        {{ project.title }}
                      </h4>
                    </div>
                    <p class="card-excerpt">{{ getPlainText(project.content) || '展示我的项目作品...' }}</p>
                  </div>
                  <div class="card-thumbnail">
                    <img :src="project.image" :alt="project.title" />
                  </div>
                  <div class="card-meta">
                    <div class="card-date">
                      <i class="date-icon">🕐</i>
                      <span>{{ formatDateTime(project.CreatedAt) }}</span>
                    </div>
                    <div class="card-tags">
                      <i class="tag-icon">🔖</i>
                      <template v-if="project.tags && project.tags.length > 0">
                        <span v-for="tag in project.tags" :key="tag" class="tag">{{ tag }}</span>
                      </template>
                      <span v-else class="tag">山角函兽懒得加标签了</span>
                    </div>
                  </div>
                </div>
              </article>
            </div>
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="sidebar">
          <!-- 天气卡片 -->
          <div class="weather-card">
            <div class="weather-header">
              <div class="location-info">
                <i class="location-icon">📍</i>
                <span class="location">{{ weatherInfo.location }}</span>
              </div>
              <div class="weather-icon">{{ getWeatherIcon(weatherInfo.weather || '') }}</div>
            </div>

            <div class="weather-main">
              <div class="temperature-section">
                <div class="temperature">{{ weatherInfo.temperature }}</div>
                <div class="weather-desc">{{ weatherInfo.weather }}</div>
              </div>
            </div>

            <div class="weather-details">
              <div v-if="weatherInfo.tomorrow" class="detail-row">
                <div class="detail-item">
                  <i class="detail-icon">🌅</i>
                  <div class="detail-content">
                    <span class="detail-label">明天</span>
                    <span class="detail-value">{{ weatherInfo.tomorrow }}</span>
                  </div>
                </div>
              </div>

              <div v-if="weatherInfo.lifeIndex" class="detail-row">
                <div class="detail-item">
                  <i class="detail-icon">💡</i>
                  <div class="detail-content">
                    <span class="detail-label">生活指数</span>
                    <span class="detail-value">{{ weatherInfo.lifeIndex }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="weather-footer">
              <div class="update-time">实时更新</div>
            </div>
          </div>

          <!-- 置顶博客 -->
          <div class="sidebar-section">
            <h4 class="section-title">置顶博客</h4>
            <div class="top-blogs">
              <div
                v-for="article in topArticles"
                :key="article.ID"
                class="top-blog-item"
                @click="$router.push(`/${article.type}/${article.ID}`)"
              >
                <div class="blog-title">
                  <font-awesome-icon :icon="getTypeIcon(article.type)" class="type-icon" />
                  {{ article.title }}
                </div>
                <div class="blog-date">{{ new Date(article.CreatedAt).toLocaleDateString('zh-CN') }}</div>
              </div>
            </div>
          </div>

          <!-- 常用链接 -->
          <div class="sidebar-section">
            <h4 class="section-title">常用链接</h4>
            <div class="useful-links">
              <div class="nav-link-item" @click="openExternal('https://github.com/Keranthos')">
                <div class="nav-link-icon" style="background: rgba(0,0,0,0.08);">
                  <font-awesome-icon :icon="['fab','github']" />
                </div>
                <div class="nav-link-content">
                  <h5>Keranthos 的 GitHub</h5>
                </div>
              </div>
            </div>
          </div>

          <!-- 浮动按钮已删除 -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import NavBar from '@/components/NavBar'
import ModernLoading from '@/components/ModernLoading.vue'
import { getArticlesList, getTopArticles, getArticlesNum } from '@/api/Articles/browse'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const typingText = '欢迎来到山角函兽的小窝'
const displayedText = ref('')
const openExternal = (url) => {
  window.open(url, '_blank')
}
const latestBlogs = ref([])
const latestMoments = ref([])
const latestProjects = ref([])
const blogCount = ref(0)
const momentCount = ref(0)
const projectCount = ref(0)
const topArticles = ref([])
const weatherInfo = ref({
  location: '获取中...',
  temperature: '--',
  weather: '获取中...',
  tomorrow: '',
  lifeIndex: ''
})

// 确保组件正确初始化
const isComponentReady = ref(false)
const loadingProgress = ref(0)

const typeText = async () => {
  for (let i = 0; i < typingText.length; i++) {
    displayedText.value += typingText[i]
    await new Promise(resolve => setTimeout(resolve, 100)) // 每个字母延迟100ms
  }
}

const loadLatestData = async () => {
  try {
    const [blogs, moments, projects, blogNum, momentNum, projectNum] = await Promise.all([
      getArticlesList('blog', 1, 3),
      getArticlesList('moment', 1, 3),
      getArticlesList('project', 1, 3),
      getArticlesNum('blog'),
      getArticlesNum('moment'),
      getArticlesNum('project')
    ])
    latestBlogs.value = blogs.data
    latestMoments.value = moments.data
    latestProjects.value = projects.data
    blogCount.value = (blogNum && (blogNum.total || blogNum.count || blogNum.num)) ? (blogNum.total || blogNum.count || blogNum.num) : (blogs.total || blogs.count || blogs.data?.length || 0)
    momentCount.value = (momentNum && (momentNum.total || momentNum.count || momentNum.num)) ? (momentNum.total || momentNum.count || momentNum.num) : (moments.total || moments.count || moments.data?.length || 0)
    projectCount.value = (projectNum && (projectNum.total || projectNum.count || projectNum.num)) ? (projectNum.total || projectNum.count || projectNum.num) : (projects.total || projects.count || projects.data?.length || 0)
  } catch (error) {
    console.error('Failed to load latest data:', error)
  }
}

// 获取置顶文章
// 天气图标映射
const getWeatherIcon = (weather) => {
  const iconMap = {
    晴: '☀️',
    多云: '⛅',
    阴: '☁️',
    雨: '🌧️',
    雪: '❄️',
    雾: '🌫️',
    霾: '😷',
    雷: '⛈️',
    风: '💨'
  }

  for (const [key, icon] of Object.entries(iconMap)) {
    if (weather.includes(key)) {
      return icon
    }
  }
  return '🌤️' // 默认图标
}

// 处理Markdown格式，提取纯文本
const getPlainText = (content) => {
  if (!content) return ''

  // 移除Markdown格式符号
  const plainText = content
    .replace(/#{1,6}\s+/g, '') // 移除标题符号
    .replace(/\*\*(.*?)\*\*/g, '$1') // 移除粗体
    .replace(/\*(.*?)\*/g, '$1') // 移除斜体
    .replace(/`(.*?)`/g, '$1') // 移除行内代码
    .replace(/```[\s\S]*?```/g, '') // 移除代码块
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // 移除链接，保留文本
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '') // 移除图片
    .replace(/^\s*[-*+]\s+/gm, '') // 移除列表符号
    .replace(/^\s*\d+\.\s+/gm, '') // 移除有序列表
    .replace(/>\s*/g, '') // 移除引用符号
    .replace(/\n+/g, ' ') // 将换行符替换为空格
    .replace(/\s+/g, ' ') // 合并多个空格
    .trim()

  // 截取前100个字符
  return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText
}

// 格式化日期时间
const formatDateTime = (dateString) => {
  const date = new Date(dateString)
  const month = date.toLocaleDateString('en-US', { month: 'short' })
  const day = date.getDate().toString().padStart(2, '0')
  const year = date.getFullYear()
  const time = date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
  return `${month} ${day}, ${year} ${time}`
}

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

const loadTopArticles = async () => {
  try {
    // 调用获取置顶文章的API
    // console.log('正在获取置顶文章...')
    const response = await getTopArticles(3)
    // console.log('置顶文章API响应:', response)
    // console.log('置顶文章数量:', response.data ? response.data.length : 0)

    if (response.success && response.data && response.data.length > 0) {
      // console.log('找到置顶文章:', response.data.length, '篇')
      topArticles.value = response.data
    } else {
      // console.log('没有置顶文章，使用最新文章作为备选')
      // 如果没有置顶文章，从所有类型中获取最新文章
      const [blogs, moments, projects] = await Promise.all([
        getArticlesList('blog', 1, 10), // 获取更多文章用于排序
        getArticlesList('moment', 1, 10),
        getArticlesList('project', 1, 10)
      ])

      const allArticles = []
      // 合并所有文章并添加类型标识
      if (blogs.data && blogs.data.length > 0) {
        blogs.data.forEach(article => {
          allArticles.push({ ...article, type: 'blog' })
        })
      }
      if (projects.data && projects.data.length > 0) {
        projects.data.forEach(article => {
          allArticles.push({ ...article, type: 'project' })
        })
      }
      if (moments.data && moments.data.length > 0) {
        moments.data.forEach(article => {
          allArticles.push({ ...article, type: 'moment' })
        })
      }

      // 按创建时间排序，取最新的3篇
      allArticles.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt))
      const topList = allArticles.slice(0, 3)

      // console.log('备选文章:', topList.length, '篇')
      topArticles.value = topList
    }
  } catch (error) {
    console.error('Failed to load top articles:', error)
    // 降级到最新文章
    try {
      const [blogs, moments, projects] = await Promise.all([
        getArticlesList('blog', 1, 10),
        getArticlesList('moment', 1, 10),
        getArticlesList('project', 1, 10)
      ])

      const allArticles = []
      // 合并所有文章并添加类型标识
      if (blogs.data && blogs.data.length > 0) {
        blogs.data.forEach(article => {
          allArticles.push({ ...article, type: 'blog' })
        })
      }
      if (projects.data && projects.data.length > 0) {
        projects.data.forEach(article => {
          allArticles.push({ ...article, type: 'project' })
        })
      }
      if (moments.data && moments.data.length > 0) {
        moments.data.forEach(article => {
          allArticles.push({ ...article, type: 'moment' })
        })
      }

      // 按创建时间排序，取最新的3篇
      allArticles.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt))
      const topList = allArticles.slice(0, 3)

      // console.log('降级到最新文章:', topList.length, '篇')
      topArticles.value = topList
    } catch (fallbackError) {
      console.error('Failed to load fallback articles:', fallbackError)
    }
  }
}

// 获取天气信息
const getWeatherInfo = async () => {
  try {
    // 从后端API获取天气信息
    const response = await fetch('/api/weather')
    const data = await response.json()

    if (data.success && data.data) {
      weatherInfo.value = {
        location: data.data.location,
        temperature: data.data.temperature,
        weather: data.data.weather,
        tomorrow: data.data.tomorrow || '',
        lifeIndex: data.data.life_index || ''
      }
      // console.log('天气信息获取成功:', data.data)
    } else {
      throw new Error('天气API返回错误')
    }
  } catch (error) {
    // console.warn('天气获取失败，使用默认天气:', error)
    // 降级到默认天气
    weatherInfo.value = {
      location: '北京',
      temperature: '22°',
      weather: '多云',
      tomorrow: '晴 17°C/5°C',
      lifeIndex: '适宜出行'
    }
  }
}

onMounted(async () => {
  try {
    // 模拟加载进度
    const progressInterval = setInterval(() => {
      if (loadingProgress.value < 90) {
        loadingProgress.value += Math.random() * 15
      }
    }, 200)

    await typeText()
    loadingProgress.value = 25

    await loadLatestData()
    loadingProgress.value = 50

    await loadTopArticles()
    loadingProgress.value = 75

    await getWeatherInfo()
    loadingProgress.value = 100

    // 清除进度条更新
    clearInterval(progressInterval)

    // 延迟一点时间让用户看到100%的进度
    setTimeout(() => {
      isComponentReady.value = true
    }, 500)
  } catch (error) {
    console.error('组件初始化失败:', error)
    isComponentReady.value = true // 即使出错也标记为就绪
  }
})
</script>

<style scoped>
/* 旧的加载样式已移除，现在使用 ModernLoading 组件 */

.home-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  overflow-y: auto;
}

/* Hero Section */
.hero-section {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  width: 100%;
  padding: 0 20px;
}

.hero-content {
  text-align: center;
  z-index: 2;
}

.avatar-container {
  position: relative;
  width: 150px;
  height: 150px;
  margin: 0 auto 30px;
}

.hero-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  transition: transform 0.3s ease;
}

.hero-avatar:hover {
  transform: scale(1.05);
}

.avatar-ring {
  position: absolute;
  top: -10px;
  left: -10px;
  right: -10px;
  bottom: -10px;
  border-radius: 50%;
  border: 2px solid rgba(106, 27, 154, 0.3);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.5;
  }
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 700;
  color: white;
  margin-bottom: 20px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  letter-spacing: 2px;
}

.cursor {
  display: inline-block;
  animation: blink 1s infinite;
  margin-left: 5px;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.hero-subtitle {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 40px;
  letter-spacing: 3px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.hero-stats {
  display: flex;
  gap: 60px;
  justify-content: center;
  margin-top: 40px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 30px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  min-width: 120px;
  flex: 0 0 auto;
}

.stat-item:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.stat-number {
  font-size: 2.5rem;
  font-weight: bold;
  color: white;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.8);
  letter-spacing: 1px;
}

.scroll-indicator {
  position: absolute;
  bottom: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.9rem;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
  60% {
    transform: translateY(-5px);
  }
}

.scroll-arrow {
  font-size: 1.5rem;
  margin-top: 10px;
}

/* 现代化内容区域 - 模仿kirigaya.cn */
.modern-content {
  width: 100%;
  padding: 40px 20px;
  background: transparent;
  min-height: 100vh;
}

.main-layout {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 0;
  align-items: start;
  position: relative;
}

/* 添加竖线分隔 */
.main-layout::after {
  content: '';
  position: absolute;
  left: calc(66.666% - 1px);
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(180deg, rgba(168, 85, 247, 0.3) 0%, rgba(124, 58, 237, 0.3) 100%);
  z-index: 1;
}

.main-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
  padding-right: 20px;
  position: relative;
  z-index: 2;
}

/* 主内容区域样式 - 与背景融为一体，无边框 */
.modern-content .main-content .content-section {
  background: transparent !important;
  backdrop-filter: none !important;
  border-radius: 0 !important;
  padding: 25px 0 !important;
  border: none !important;
  margin: 0 !important;
  min-height: auto !important;
}

/* 区域标题样式 */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 2px solid rgba(168, 85, 247, 0.1);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.4rem;
  font-weight: 700;
  color: #333;
  margin: 0;
}

.section-icon {
  font-size: 1.6rem;
}

.section-text {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-count {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

/* 博客卡片网格布局 */
.blog-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 40px;
}

/* 现代化卡片设计 - 与背景融为一体，无边框 */
.blog-card {
  background: transparent;
  border-radius: 0;
  border: none;
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.blog-card:hover {
  background: linear-gradient(135deg, rgba(248, 249, 255, 0.1) 0%, rgba(240, 242, 255, 0.1) 100%);
  transform: translateY(-2px);
}

/* 卡片内容布局 - 左右分布，65%文字 + 35%图片 */
.card-content {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 20px;
  min-height: 100px;
  border-radius: 15px;
  transition: all 0.3s ease;
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

.card-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-icon {
  font-size: 1rem;
  margin-right: 8px;
  color: #a855f7;
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
  line-height: 1.6;
  margin: 0;
  text-align: left;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 卡片元信息 - 底部左右分布 */
.card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 15px;
  border-top: 1px solid rgba(168, 85, 247, 0.1);
  width: 100%;
  /* 让分割线铺满卡片内容宽度（抵消内边距）*/
  margin-left: -20px;
  margin-right: -20px;
  padding-left: 20px;
  padding-right: 20px;
}

.card-date {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #999;
  font-size: 0.85rem;
  font-weight: 500;
}

.date-icon {
  font-size: 0.9rem;
}

.card-tags {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap; /* 与图标同一行，不换行 */
  overflow: hidden;  /* 超出隐藏，配合单行显示 */
  min-width: 0;      /* 允许在 flex 布局中收缩并触发省略号 */
  margin-left: auto; /* 容器整体靠右，但内容从左开始，便于右侧出现 ... */
}

.tag-icon {
  font-size: 0.9rem;
  color: #999;
}

.tag {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 500;
  border: 1px solid rgba(168, 85, 247, 0.2);
  transition: all 0.3s ease;
}

.tag:hover {
  background: rgba(168, 85, 247, 0.2);
  transform: translateY(-1px);
}

/* 卡片缩略图 - 右侧，35%宽度，高度增加 */
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

/* 侧边栏样式 */
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 25px;
  position: sticky;
  top: 100px;
  padding-left: 20px;
  position: relative;
  z-index: 2;
}

.weather-card {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  border-radius: 20px;
  padding: 24px;
  color: white;
  box-shadow: 0 12px 40px rgba(168, 85, 247, 0.25);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.weather-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  animation: shimmer 3s ease-in-out infinite;
}

@keyframes shimmer {
  0%, 100% { transform: translateX(-100%) translateY(-100%) rotate(0deg); }
  50% { transform: translateX(0%) translateY(0%) rotate(180deg); }
}

.weather-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  position: relative;
  z-index: 2;
}

.location-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.location-icon {
  font-size: 16px;
  opacity: 0.9;
}

.location {
  font-size: 16px;
  font-weight: 600;
  opacity: 0.95;
}

.weather-icon {
  font-size: 32px;
  opacity: 0.9;
  animation: float 2s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-5px); }
}

.location {
  font-size: 0.9rem;
  opacity: 0.9;
}

.weather-main {
  margin-bottom: 24px;
  position: relative;
  z-index: 2;
}

.temperature-section {
  text-align: center;
}

.temperature {
  font-size: 3.5rem;
  font-weight: 800;
  line-height: 1;
  margin-bottom: 8px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  background: linear-gradient(45deg, #fff, #f0f8ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.weather-desc {
  font-size: 1.2rem;
  opacity: 0.9;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.weather-desc {
  font-size: 1.1rem;
  font-weight: 500;
}

.weather-details {
  margin-bottom: 20px;
  position: relative;
  z-index: 2;
}

.detail-row {
  margin-bottom: 16px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.detail-item:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.detail-icon {
  font-size: 20px;
  opacity: 0.9;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.detail-label {
  font-size: 0.85rem;
  opacity: 0.8;
  font-weight: 500;
}

.detail-value {
  font-size: 0.95rem;
  font-weight: 600;
  opacity: 0.95;
}

.weather-footer {
  text-align: center;
  position: relative;
  z-index: 2;
}

.update-time {
  font-size: 0.8rem;
  opacity: 0.7;
  font-weight: 400;
  letter-spacing: 0.5px;
}

.sidebar-section {
  background: transparent;
  backdrop-filter: none;
  border-radius: 0;
  padding: 20px 0;
  box-shadow: none;
  border: none;
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 15px;
  padding-bottom: 8px;
  border-bottom: 2px solid rgba(168, 85, 247, 0.1);
}

/* 导航链接样式 */
.nav-links {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.nav-link-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.nav-link-item:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.nav-link-icon {
  font-size: 1.5rem;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.nav-link-content h5 {
  margin: 0 0 4px 0;
  font-size: 1rem;
  font-weight: 600;
  color: #333;
}

.nav-link-content p {
  margin: 0;
  font-size: 0.85rem;
  color: #666;
}

.top-blogs {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.top-blog-item {
  padding: 10px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.top-blog-item:last-child {
  border-bottom: none;
}

.top-blog-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(168, 85, 247, 0.1), transparent);
  transition: left 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.top-blog-item:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow:
    0 8px 25px rgba(168, 85, 247, 0.15),
    0 4px 12px rgba(0, 0, 0, 0.1);
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.08) 0%, rgba(124, 58, 237, 0.08) 100%);
  border-radius: 8px;
  padding: 12px 16px;
  margin: 0 -16px 4px -16px;
  border: 1px solid rgba(168, 85, 247, 0.2);
  border-bottom: 1px solid rgba(168, 85, 247, 0.2);
}

.top-blog-item:hover::before {
  left: 100%;
}

.blog-title {
  font-size: 0.9rem;
  color: #333;
  margin-bottom: 6px;
  line-height: 1.4;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  z-index: 1;
}

.blog-title:hover {
  color: #a855f7;
  transform: translateX(4px);
}

.blog-date {
  font-size: 0.8rem;
  color: #999;
  text-align: left;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.top-blog-item:hover .blog-date {
  color: #a855f7;
  transform: translateX(4px);
}

.type-icon {
  font-size: 0.9rem;
  opacity: 0.7;
  flex-shrink: 0;
  color: #a855f7;
  margin-right: 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 1;
}

.top-blog-item:hover .type-icon {
  opacity: 1;
  transform: scale(1.1) rotate(5deg);
  color: #5a67d8;
}

.useful-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  cursor: pointer;
  transition: color 0.2s ease;
}

.link-item:hover {
  color: #a855f7;
}

.link-icon {
  font-size: 1.2rem;
}

/* 浮动按钮样式已删除 */

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-layout {
    grid-template-columns: 1fr;
    gap: 30px;
  }

  .sidebar {
    position: static;
    order: -1;
  }

  /* 浮动按钮已删除 */
}

/* 当窗口 <= 1330px（或相当于最大布局宽度的 2/3）时，隐藏右侧栏与分隔线，仅显示左侧内容 */
@media (max-width: 1330px) {
  .main-layout {
    grid-template-columns: 1fr;
  }
  .main-layout::after { display: none; }
  .sidebar { display: none; }
  .main-content {
    padding-right: 0;
    width: 66.666%;
    margin: 0 auto;
  }
  /* 紧凑模式下隐藏首页统计块，并避免与固定导航重叠 */
  .hero-stats { display: none; }
  .hero-section { margin-top: 40px; }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .hero-stats {
    flex-direction: column;
    gap: 20px;
  }

  .stat-item {
    min-width: 100px;
    padding: 15px 20px;
  }

  .modern-content {
    padding: 20px 15px;
  }

  .main-layout {
    gap: 20px;
  }
  /* 手机端也保持主体 2/3 宽度，左右留边距 */
  .main-content {
    width: 66.666%;
    margin: 0 auto;
    min-width: 480px; /* 手机端主体区域最小宽度不小于 480px */
  }

  .blog-cards {
    gap: 15px;
  }

  .card-content {
    gap: 12px;
    padding: 15px;
    min-height: 80px;
    border-radius: 12px;
    display: flex;
    flex-direction: column; /* 垂直布局，便于重排 */
  }

  .card-text { order: 1; }

  .card-thumbnail {
    order: 2; /* 图片位于简介与底部行之间 */
    width: 100%;
    height: auto;
    /* 不超过卡片高度的 2/5（若父级有高度时生效）；并提供基于视口宽度的兜底限制 */
    max-height: 40%;
    max-height: min(40vw, 200px);
    border-radius: 10px;
  }
  .card-thumbnail img { width: 100%; height: 100%; object-fit: cover; max-height: inherit; }

  .card-title {
    font-size: 1.1rem;
  }

  .card-meta {
    order: 3; /* 底部行 */
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-top: 10px;
    padding-top: 10px;
    /* 移动端内边距为15px，分割线同样需要满宽 */
    margin-left: -15px;
    margin-right: -15px;
    padding-left: 15px;
    padding-right: 15px;
  }
  /* 标签容器固定占 2/5 宽，靠右并保持与图标同一行 */
  .card-tags {
    flex: 0 0 40%;
    /* 容器靠右，由自身位置保证贴右，内容左起以便右侧出现 ... */
    margin-left: auto;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .card-tags .tag { display: inline; padding: 0 6px; }

  .sidebar-section {
    padding: 15px;
  }

  .weather-card {
    padding: 20px;
  }

  .temperature {
    font-size: 2.8rem;
  }

  .weather-icon {
    font-size: 28px;
  }

  .detail-item {
    padding: 10px 12px;
  }
}

@media (max-width: 480px) {
  .section-title {
    font-size: 1.2rem;
  }

  .section-count {
    font-size: 0.8rem;
    padding: 4px 8px;
  }

  .blog-cards {
    gap: 12px;
  }

  .card-content {
    padding: 12px;
    gap: 10px;
  }

  .card-thumbnail {
    height: 140px; /* 略大一些，位于中间更协调 */
  }

  .card-title {
    font-size: 1rem;
  }

  .card-subtitle {
    font-size: 0.9rem;
  }

  .card-excerpt {
    font-size: 0.85rem;
  }

  .card-meta {
    margin-top: 8px;
    padding-top: 8px;
  }

  .card-date {
    font-size: 0.8rem;
  }

  .tag {
    font-size: 0.75rem;
    padding: 3px 6px;
  }

  /* 浮动按钮已删除 */
}
</style>
