<template>
  <div class="article-list-view">
    <NavBar />
    <!-- 头部图片区域 -->
    <div class="header-section">
      <img :src="headerImage" alt="Header Image" class="header-image" />
    </div>

    <!-- 标签页区域 -->
    <div class="tabs-section">
      <div class="tabs-container">
        <div class="tab-background" :style="tabBackgroundTransform"></div>
        <div class="tab-item" :class="{ active: activeTab === 'main' }" @click="switchTab('main')">
          <span class="tab-icon">📚</span>
          <span class="tab-text">{{ getTypeName(type) }}</span>
          <span class="tab-count">({{ getArticleCount() }})</span>
        </div>
        <div class="tab-item" :class="{ active: activeTab === 'thoughts' }" @click="switchTab('thoughts')">
          <span class="tab-icon">🌟</span>
          <span class="tab-text">所思所想</span>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="content-section">
      <!-- 主要内容（博客/项目/科研） -->
      <div v-if="activeTab === 'main'">
        <div v-if="loading" class="loading-wrapper">
          <ModernLoading
            :progress="loadingProgress"
            :title="getTypeName(type)"
            :subtitle="'Loading……'"
          />
        </div>
        <div v-else-if="articles.length > 0" class="article-grid">
          <ArticleCard
            v-for="article in articles"
            :id="article.ID" :key="article.ID" :image="article.image"
            :title="article.title" :tags="article.tags" :time="article.CreatedAt"
            :type="type" :reading-time="article.readingTime" :article-type="article.articleType"
          />
        </div>
        <div v-else class="no-content">
          <p>暂无内容</p>
        </div>
      </div>

      <!-- 所思所想内容 -->
      <div v-if="activeTab === 'thoughts'" class="thoughts-content">
        <div class="thoughts-container">
          <h2 class="thoughts-title">{{ thoughtsContent.title }}</h2>
          <div class="thoughts-text">
            <p v-for="(paragraph, index) in thoughtsContent.paragraphs" :key="index" class="thoughts-paragraph">
              {{ paragraph }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页区域 -->
    <div v-if="activeTab === 'main'" class="pagination-section">
      <div class="pagination">
        <button :disabled="currentPage <= 1" @click="loadPage(currentPage - 1)">
          上一页
        </button>
        <button v-if="currentPage > 4" @click="loadPage(1)">1</button>
        <span v-if="currentPage > 4">...</span>
        <button v-for="page in pagesToShow" :key="page" :disabled="page === currentPage" @click="loadPage(page)">
          {{ page }}
        </button>
        <span v-if="currentPage < totalPage - 3">...</span>
        <button v-if="currentPage < totalPage - 3" @click="loadPage(totalPage)">{{ totalPage }}</button>
        <button :disabled="currentPage >= totalPage" @click="loadPage(currentPage + 1)">
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// 上条栏中博客、项目、科研日记的基础显示视图
import { ref, onMounted, watch, computed, defineAsyncComponent } from 'vue'
// import ArticleCard from '@/components/ArticleCard.vue'
import { getArticlesList, getArticlesNum } from '@/api/Articles/browse'
import { batchEstimateReadingTime } from '@/utils/readingTime'
// import { NavBar } from '@/components/NavBar'
const ArticleCard = defineAsyncComponent(() => import('@/components/ArticleCard.vue'))
const NavBar = defineAsyncComponent(() => import('@/components/NavBar.vue'))
const ModernLoading = defineAsyncComponent(() => import('@/components/ModernLoading.vue'))

const props = defineProps({
  type: String, // 'blog', 'project', 'research', 'moment'
  headerImage: String,
  typingText: String
})

const currentPage = ref(1)
const totalPage = ref(1)
const limit = 9
const articles = ref([])
const articleCache = ref({})// 用于缓存数据
const displayedText = ref('')
const loading = ref(false)
const loadingProgress = ref(0)
const activeTab = ref('main') // 当前激活的标签页
const totalArticles = ref(0) // 实际文章总数

// 所思所想内容
const thoughtsContent = ref({
  title: '',
  paragraphs: []
})

// 加载所思所想内容
const loadThoughtsContent = async () => {
  try {
    const response = await fetch('/thoughts.txt')
    const text = await response.text()

    // 解析文本内容
    const lines = text.split('\n').filter(line => line.trim() !== '')
    if (lines.length > 0) {
      thoughtsContent.value.title = lines[0]
      thoughtsContent.value.paragraphs = lines.slice(1)
    }
  } catch (error) {
    console.error('加载所思所想内容失败:', error)
    // 设置默认内容
    thoughtsContent.value = {
      title: '时间为何只是单向地流淌?',
      paragraphs: [
        '我一直在思考，为什么要写博客。最开始的时候，我在知乎上写一些系统性的文章，在CSDN上写一些技术笔记。后来，我觉得应该有一个属于自己的地方，可以记录自己的思考，可以分享自己的经验，可以与他人交流。',
        '于是，我开始在知乎上创建专栏，开始写技术博客，开始分享自己的学习心得。我希望通过写作，能够帮助更多的人，也能够让自己在技术道路上不断成长。同时，我也希望能够在互联网的世界里，留下一些有价值的内容，让知识得以传承。',
        '时间为何只是单向地流淌？我想，也许是因为我们总是想要抓住些什么，留住些什么。而写作，就是我在时间的长河中，留下的一些痕迹。这些痕迹，或许能够帮助到一些人，或许能够启发一些思考，或许只是我对自己存在的一种证明。'
      ]
    }
  }
}

const fetchArticlesNum = async () => {
  try {
    // console.log('fetchArticlesNum called with type:', props.type)
    if (props.type === 'moment') {
      const { getMomentsNum } = await import('@/api/Moments/browse')
      const response = await getMomentsNum()
      // console.log('Moments response:', response)
      totalArticles.value = response.num
      return response.num
    } else if (props.type === 'all') {
      // 获取所有类型文章的总数
      // console.log('Fetching all article types...')
      const blogResponse = await getArticlesNum('blog')
      const projectResponse = await getArticlesNum('project')
      const researchResponse = await getArticlesNum('research')
      // console.log('Blog count:', blogResponse.num)
      // console.log('Project count:', projectResponse.num)
      // console.log('Research count:', researchResponse.num)
      const total = blogResponse.num + projectResponse.num + researchResponse.num
      // console.log('Total articles:', total)
      totalArticles.value = total
      return total
    } else {
      const response = await getArticlesNum(props.type)
      // console.log(`${props.type} response:`, response)
      totalArticles.value = response.num
      return response.num
    }
  } catch (error) {
    console.error('Error in fetchArticlesNum:', error)
    totalArticles.value = 0
    return 0
  }
}

const loadPage = async (page) => {
  if (page < 1 || page > totalPage.value) return
  currentPage.value = page

  // 检查缓存中是否已有该页的数据
  const cacheKey = `${props.type}-${page}`
  if (articleCache.value[cacheKey]) {
    articles.value = articleCache.value[cacheKey]
  } else {
    loading.value = true
    loadingProgress.value = 0

    // 模拟加载进度
    const progressInterval = setInterval(() => {
      if (loadingProgress.value < 90) {
        loadingProgress.value += Math.random() * 20
      }
    }, 100)

    try {
      if (props.type === 'moment') {
        const { getMomentsList } = await import('@/api/Moments/browse')
        const response = await getMomentsList(page, limit)
        // 映射字段以匹配 ArticleCard 组件
        const momentArticles = response.data.map(item => ({
          ID: item.ID,
          image: item.Image || 'https://picsum.photos/id/201/800/600',
          title: item.Title,
          content: item.Content,
          tags: [],
          CreatedAt: item.CreatedAt
        }))

        // 批量计算阅读时间
        articles.value = batchEstimateReadingTime(momentArticles)
        articleCache.value[cacheKey] = articles.value
      } else if (props.type === 'all') {
        // 获取所有类型的文章（博客、项目、科研）
        const blogResponse = await getArticlesList('blog', page, limit)
        const projectResponse = await getArticlesList('project', page, limit)
        const researchResponse = await getArticlesList('research', page, limit)

        // 合并所有文章并添加类型标识
        const allArticles = [
          ...blogResponse.data.map(item => ({ ...item, articleType: 'blog' })),
          ...projectResponse.data.map(item => ({ ...item, articleType: 'project' })),
          ...researchResponse.data.map(item => ({ ...item, articleType: 'research' }))
        ]

        // 按创建时间排序
        allArticles.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt))

        // 批量计算阅读时间
        articles.value = batchEstimateReadingTime(allArticles)
        articleCache.value[cacheKey] = articles.value
      } else {
        const response = await getArticlesList(props.type, page, limit)

        // 批量计算阅读时间
        articles.value = batchEstimateReadingTime(response.data)
        articleCache.value[cacheKey] = articles.value // 缓存该页的数据
      }

      // 完成加载
      loadingProgress.value = 100
      clearInterval(progressInterval)

      // 延迟一点时间让用户看到100%的进度
      setTimeout(() => {
        loading.value = false
      }, 300)
    } catch (error) {
      console.error('加载页面失败:', error)
      clearInterval(progressInterval)
      loading.value = false
    }
  }
}

const pagesToShow = computed(() => {
  const pages = []
  const startPage = Math.max(1, currentPage.value - 3)
  const endPage = Math.min(totalPage.value, currentPage.value + 3)

  for (let i = startPage; i <= endPage; i++) {
    pages.push(i)
  }

  return pages
})

// 获取类型名称
const getTypeName = (type) => {
  const typeMap = {
    blog: '我的博客',
    project: '我的项目',
    research: '我的科研',
    moment: '碎碎念',
    all: '我的博客'
  }
  return typeMap[type] || type
}

// 获取文章数量
const getArticleCount = () => {
  return totalArticles.value
}

// 切换标签页
const switchTab = (tab) => {
  activeTab.value = tab
}

// 计算背景位置和宽度
const tabBackgroundTransform = computed(() => {
  if (activeTab.value === 'main') {
    return {
      transform: 'translateX(0)',
      width: '160px'
    }
  } else {
    // 对于第二个标签，需要计算第一个标签的宽度
    // 由于标签宽度是自适应的，我们使用一个估算值
    return {
      transform: 'translateX(160px)',
      width: '120px'
    }
  }
})

onMounted(async () => {
  displayedText.value = props.typingText
  await loadThoughtsContent() // 加载所思所想内容
  const articleNum = await fetchArticlesNum()
  totalPage.value = Math.ceil(articleNum / limit)
  loadPage(1)
})

watch(() => props.type, async () => {
  const articleNum = await fetchArticlesNum()
  totalPage.value = Math.ceil(articleNum / limit)
  loadPage(1)
})
</script>

<style scoped>
.article-list-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: transparent;
}

/* 头部图片区域 */
.header-section {
  position: relative;
  width: 100%;
  margin-top: 80px;
  padding: 40px 350px;
  overflow: hidden;
}

.header-image {
  width: 100%;
  height: 320px;
  object-fit: cover;
  border-radius: 8px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  transition: transform 0.6s ease;
  position: relative;
  z-index: 1;
}

.header-image:hover {
  transform: scale(1.05);
}

/* 标签页区域 */
.tabs-section {
  width: 100%;
  padding: 0 350px;
  margin-top: 10px;
}

.tabs-container {
  display: flex;
  gap: 0;
  background: transparent;
  border-bottom: 3px solid #6a1b9a;
  padding: 0;
  position: relative;
}

.tab-background {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(135deg, rgba(106, 27, 154, 0.9) 0%, rgba(106, 27, 154, 0.7) 100%);
  border-radius: 12px 12px 0 0;
  border-bottom: none;
  transition: transform 0.3s ease, width 0.3s ease;
  z-index: 1;
}

.tab-item {
  padding: 10px 16px;
  background: transparent;
  border: none;
  border-radius: 12px 12px 0 0;
  cursor: pointer;
  transition: none;
  position: relative;
  font-size: 0.95rem;
  font-weight: 500;
  color: #333;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: none;
  border-bottom: 3px solid transparent;
  z-index: 2;
  width: max-content;
  min-width: fit-content;
}

.tab-item .tab-icon {
  filter: brightness(0);
}

.tab-item.active {
  color: white;
}

.tab-item.active .tab-icon {
  filter: brightness(0) invert(1);
}

.tab-item:hover {
  cursor: pointer;
}

.tab-text {
  margin-right: 8px;
}

.tab-count {
  font-size: 0.9rem;
  opacity: 0.8;
  font-weight: 500;
}

.tab-icon {
  font-size: 1.1rem;
}

/* 所思所想内容样式 - 与背景融为一体 */
.thoughts-content {
  width: 100%;
  background: transparent; /* 移除白色背景 */
  border-radius: 0; /* 移除圆角 */
  box-shadow: none; /* 移除阴影 */
  overflow: visible; /* 允许内容溢出 */
}

.thoughts-container {
  padding: 0; /* 移除内边距 */
  text-align: left; /* 容器内文本左对齐 */
}

.thoughts-title {
  font-size: 2.2rem;
  font-weight: 700;
  color: #212529; /* 深色文字，在背景上更清晰 */
  margin-bottom: 40px;
  line-height: 1.3;
  text-align: left; /* 标题左对齐 */
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 添加文字阴影增强可读性 */
}

.thoughts-text {
  line-height: 1.8;
  text-align: left; /* 文本内容左对齐 */
}

.thoughts-paragraph {
  font-size: 1.2rem;
  color: #495057; /* 深色文字，在背景上更清晰 */
  margin-bottom: 24px;
  text-indent: 2em;
  text-align: left; /* 段落左对齐 */
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1); /* 添加文字阴影增强可读性 */
}

.thoughts-paragraph:last-child {
  margin-bottom: 0;
}

/* 内容区域 */
.content-section {
  width: 100%;
  padding: 40px 350px;
  background: transparent;
  position: relative;
}
.article-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  row-gap: 80px;
  column-gap: 40px;
  margin-bottom: 40px;
  position: relative;
}
/* 分页区域 */
.pagination-section {
  width: 100%;
  padding: 40px 350px;
  display: flex;
  justify-content: center;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.pagination button {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid rgba(106, 27, 154, 0.3);
  border-radius: 10px;
  color: #6a1b9a;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.pagination button:hover:not(:disabled) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(106, 27, 154, 0.3);
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.5);
}

.pagination span {
  color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 1600px) {
  .header-section,
  .tabs-section,
  .content-section,
  .pagination-section {
    padding-left: 200px;
    padding-right: 200px;
  }
}

@media (max-width: 1200px) {
  .header-section,
  .tabs-section,
  .content-section,
  .pagination-section {
    padding-left: 100px;
    padding-right: 100px;
  }

  .article-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* 加载包装器样式 */
.loading-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
}

@media (max-width: 768px) {
  .header-section,
  .tabs-section,
  .content-section,
  .pagination-section {
    padding-left: 20px;
    padding-right: 20px;
  }

  .header-image {
    height: 250px;
  }

  .typing-text {
    font-size: 1.8rem;
    padding: 20px;
  }

  .article-grid {
    grid-template-columns: 1fr;
  }

  .tabs-container {
    flex-wrap: wrap;
  }

  .tab-item {
    padding: 12px 20px;
    font-size: 1rem;
  }

  .pagination button {
    padding: 8px 15px;
    font-size: 0.9rem;
  }
}
</style>
