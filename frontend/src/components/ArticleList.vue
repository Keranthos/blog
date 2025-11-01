<template>
  <div class="article-list-view">
    <NavBar />
    <!-- 头部图片区域 -->
    <div class="header-section">
      <img :src="headerImage" alt="Header Image" class="header-image" loading="lazy" decoding="async" @error="onImgError($event)" />
    </div>

    <!-- 标签页区域 -->
    <div class="tabs-section">
      <div ref="tabsContainer" class="tabs-container">
        <div class="tab-background" :style="tabBackgroundTransform"></div>
        <div ref="tabMain" class="tab-item" :class="{ active: activeTab === 'main' }" @click="switchTab('main')">
          <span class="tab-icon">📚</span>
          <span class="tab-text">{{ getTypeName(type) }}</span>
          <span class="tab-count">({{ getArticleCount() }})</span>
        </div>
        <div ref="tabThoughts" class="tab-item" :class="{ active: activeTab === 'thoughts' }" @click="switchTab('thoughts')">
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
        <!-- 无限滚动哨兵（仅手机端且未到最后一页时显示） -->
        <div v-if="isMobile && currentPage < totalPage" ref="infiniteSentinel" class="infinite-sentinel"></div>
        <div v-if="!isMobile && currentPage === 1 && (!articles || articles.length === 0)" class="no-content">
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
    <div v-if="activeTab === 'main' && totalPage > 1 && !isMobile" class="pagination-section">
      <div class="pagination">
        <button class="nav-btn prev" :disabled="currentPage <= 1" @click="loadPage(currentPage - 1)">
          PREV
        </button>
        <button v-if="currentPage > 4" @click="loadPage(1)">1</button>
        <span v-if="currentPage > 4">...</span>
        <button v-for="page in pagesToShow" :key="page" :disabled="page === currentPage" @click="loadPage(page)">
          {{ page }}
        </button>
        <span v-if="currentPage < totalPage - 3">...</span>
        <button v-if="currentPage < totalPage - 3" @click="loadPage(totalPage)">{{ totalPage }}</button>
        <button class="nav-btn next" :disabled="currentPage >= totalPage" @click="loadPage(currentPage + 1)">
          NEXT
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// 上条栏中博客、项目、科研日记的基础显示视图
import { ref, onMounted, onBeforeUnmount, watch, computed, nextTick, defineAsyncComponent } from 'vue'
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
const fallbackImg = '/images/sunset-mountains.jpg'

// 图片错误回退
const onImgError = (e) => {
  const img = e?.target
  if (img && img.src !== fallbackImg) {
    img.src = fallbackImg
  }
}

// 所思所想内容
const thoughtsContent = ref({
  title: '',
  paragraphs: []
})

// 加载所思所想内容
const loadThoughtsContent = async () => {
  try {
    // 根据类型加载不同的文件：blog 使用 blog-thoughts.txt，moment 使用 moment-thoughts.txt
    const fileName = props.type === 'moment' ? 'moment-thoughts.txt' : 'blog-thoughts.txt'
    const response = await fetch(`/data/${fileName}`)
    const text = await response.text()

    // 解析文本内容
    const lines = text.split('\n').filter(line => line.trim() !== '')
    if (lines.length > 0) {
      thoughtsContent.value.title = lines[0]
      thoughtsContent.value.paragraphs = lines.slice(1)
    }
  } catch (error) {
    console.error('加载所思所想内容失败:', error)
    // 设置默认内容（根据类型）
    if (props.type === 'moment') {
      thoughtsContent.value = {
        title: '关于碎碎念',
        paragraphs: [
          '生活中有很多转瞬即逝的想法，它们可能不够系统，也不够深刻，但它们真实地记录了我某个时刻的感受。',
          '有些想法像流星一样划过脑海，如果不及时记录下来，很快就会消失在记忆的海洋中。所以我把它们写在这里，不是为了给别人看，而是为了给未来的自己看。',
          '这些碎碎的念头，是我与世界的对话，也是我与自己的和解。'
        ]
      }
    } else {
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

// 是否为手机端（≤768px）
const isMobile = ref(typeof window !== 'undefined' ? window.innerWidth <= 768 : false)
const handleResize = () => { isMobile.value = window.innerWidth <= 768 }

// 无限滚动哨兵
const infiniteSentinel = ref(null)
let io = null
let loadDebounceTimer = null // 去抖定时器
const isLoadingPage = ref(false) // 加载锁，防止并发加载
const retryCount = ref(0) // 重试计数器
const maxRetries = 3 // 最大重试次数
const debounceDelay = 300 // 去抖延迟（毫秒）

// 去抖函数
const debouncedLoadPage = (page, append = false) => {
  clearTimeout(loadDebounceTimer)
  loadDebounceTimer = setTimeout(() => {
    loadPageWithRetry(page, append)
  }, debounceDelay)
}

// 带重试的加载函数
const loadPageWithRetry = async (page, append = false, retry = 0) => {
  if (isLoadingPage.value || loading.value) return // 如果正在加载，直接返回
  if (page < 1 || page > totalPage.value) return

  isLoadingPage.value = true
  retryCount.value = retry

  try {
    await loadPage(page, append)
    retryCount.value = 0 // 成功后重置重试计数
  } catch (error) {
    console.error('加载页面失败:', error)
    // 如果还有重试机会，延迟后重试
    if (retry < maxRetries) {
      const retryDelay = (retry + 1) * 1000 // 递增延迟：1s, 2s, 3s
      setTimeout(() => {
        loadPageWithRetry(page, append, retry + 1)
      }, retryDelay)
    } else {
      // 重试次数用尽，显示错误提示
      console.error('加载失败，已达到最大重试次数')
      retryCount.value = 0
    }
  } finally {
    isLoadingPage.value = false
  }
}

const setupInfiniteScroll = async () => {
  if (!isMobile.value) { teardownInfiniteScroll(); return }
  if (!('IntersectionObserver' in window)) return
  // 等待 DOM 渲染出哨兵元素
  await nextTick()
  if (!infiniteSentinel.value) return
  teardownInfiniteScroll()
  io = new IntersectionObserver((entries) => {
    const e = entries[0]
    if (e && e.isIntersecting && !loading.value && !isLoadingPage.value && currentPage.value < totalPage.value) {
      // 使用去抖函数，避免快速滚动时触发多次加载
      debouncedLoadPage(currentPage.value + 1, true)
    }
  }, { rootMargin: '200px 0px', threshold: 0.1 })
  io.observe(infiniteSentinel.value)
}

const teardownInfiniteScroll = () => {
  if (io) {
    io.disconnect()
    io = null
  }
  // 清理去抖定时器
  if (loadDebounceTimer) {
    clearTimeout(loadDebounceTimer)
    loadDebounceTimer = null
  }
}

const loadPage = async (page, append = false) => {
  if (page < 1 || page > totalPage.value) {
    return Promise.resolve() // 确保总是返回 Promise
  }
  currentPage.value = page

  // 检查缓存中是否已有该页的数据
  const cacheKey = `${props.type}-${page}`
  if (articleCache.value[cacheKey]) {
    if (append) {
      articles.value = [...articles.value, ...articleCache.value[cacheKey]]
    } else {
      articles.value = articleCache.value[cacheKey]
    }
    return Promise.resolve() // 从缓存读取成功，返回 resolved Promise
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
        const computedList = batchEstimateReadingTime(momentArticles)
        if (append) {
          articles.value = [...articles.value, ...computedList]
        } else {
          articles.value = computedList
        }
        articleCache.value[cacheKey] = computedList
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
        const computedList = batchEstimateReadingTime(allArticles)
        if (append) {
          articles.value = [...articles.value, ...computedList]
        } else {
          articles.value = computedList
        }
        articleCache.value[cacheKey] = computedList
      } else {
        const response = await getArticlesList(props.type, page, limit)

        // 批量计算阅读时间
        const computedList = batchEstimateReadingTime(response.data)
        if (append) {
          articles.value = [...articles.value, ...computedList]
        } else {
          articles.value = computedList
        }
        articleCache.value[cacheKey] = computedList // 缓存该页的数据
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
      throw error // 抛出错误，让 loadPageWithRetry 捕获并重试
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

// 标签高亮背景：基于实际 DOM 宽度与位置，避免移动端估算误差
const tabsContainer = ref(null)
const tabMain = ref(null)
const tabThoughts = ref(null)

const computeTabStyle = () => {
  const currentEl = activeTab.value === 'main' ? tabMain.value : tabThoughts.value
  const container = tabsContainer.value
  if (currentEl && container) {
    const crect = currentEl.getBoundingClientRect()
    const contRect = container.getBoundingClientRect()
    // 确保元素已完全渲染（宽度不为0）
    if (crect.width > 0) {
      const left = crect.left - contRect.left
      const width = crect.width
      return { transform: `translateX(${left}px)`, width: `${width}px` }
    }
  }
  // 回退：根据类型设置合理的初始宽度（避免闪烁）
  const estimatedWidth = props.type === 'moment' ? '140px' : (activeTab.value === 'main' ? '160px' : '120px')
  return { transform: 'translateX(0)', width: estimatedWidth }
}

// 根据类型和标签估算合理的初始宽度
const getEstimatedWidth = () => {
  if (props.type === 'moment') return '140px' // "碎碎念" 约140px
  if (activeTab.value === 'main') {
    // "我的博客" / "我的项目" / "我的科研" 约160-180px
    const typeMap = { blog: '160px', project: '160px', research: '180px', all: '160px' }
    return typeMap[props.type] || '160px'
  }
  return '120px' // "所思所想" 约120px
}

const tabBackgroundTransform = ref({ transform: 'translateX(0)', width: getEstimatedWidth() })

const updateTabBackground = async () => {
  // 先设置一个合理的估算值，避免初始显示时太小
  const estimatedWidth = getEstimatedWidth()
  tabBackgroundTransform.value = { transform: 'translateX(0)', width: estimatedWidth }

  // 等待DOM完全渲染后再精确计算
  await nextTick()
  // 使用 requestAnimationFrame 确保在浏览器下一次重绘前计算
  await new Promise(resolve => requestAnimationFrame(resolve))
  await nextTick()

  // 精确计算位置和宽度
  const computedStyle = computeTabStyle()
  // 只有当计算结果有效时才更新（避免宽度为0的情况）
  if (computedStyle.width !== '0px' && computedStyle.width !== estimatedWidth) {
    tabBackgroundTransform.value = computedStyle
  }
}

watch([activeTab, () => props.type, totalArticles], updateTabBackground)

onMounted(async () => {
  displayedText.value = props.typingText
  // 先设置初始估算值，避免显示时太小
  tabBackgroundTransform.value = { transform: 'translateX(0)', width: getEstimatedWidth() }

  await loadThoughtsContent() // 加载所思所想内容
  const articleNum = await fetchArticlesNum()
  totalPage.value = Math.ceil(articleNum / limit)
  await loadPage(1)

  // 等待所有内容渲染完成后再精确计算标签高亮位置
  // 使用多层 nextTick 和 requestAnimationFrame 确保DOM完全准备好
  await nextTick()
  await new Promise(resolve => requestAnimationFrame(resolve))
  await nextTick()
  updateTabBackground()

  // 监听窗口大小变化
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize)
    window.addEventListener('resize', updateTabBackground)
  }
  // 设置无限滚动
  setupInfiniteScroll()
})

watch(() => props.type, async () => {
  const articleNum = await fetchArticlesNum()
  totalPage.value = Math.ceil(articleNum / limit)
  // 切换类型时重置数据与缓存定位
  articles.value = []
  currentPage.value = 1
  loadPage(1)
  // 重新设置无限滚动
  teardownInfiniteScroll()
  setupInfiniteScroll()
})

watch(isMobile, () => {
  teardownInfiniteScroll()
  setupInfiniteScroll()
})

onBeforeUnmount(() => {
  teardownInfiniteScroll()
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize)
    window.removeEventListener('resize', updateTabBackground)
  }
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
  margin-top: 40px;
  padding: 40px 350px;
  overflow: hidden;
}

.header-image {
  width: 100%;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
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
  /* 设置最小宽度，确保初始显示时能包裹文字 */
  min-width: 120px;
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

.infinite-sentinel { width: 100%; height: 1px; }

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.pagination button {
  padding: 10px 18px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid rgba(106, 27, 154, 0.25);
  border-radius: 10px;
  color: #6a1b9a;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.pagination button:hover:not(:disabled) {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(106, 27, 154, 0.3);
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.5);
}

/* 更大气的上一页/下一页按钮 */
.pagination .nav-btn {
  padding: 12px 28px;
  border-radius: 14px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.95rem;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: #fff;
  border: 0;
  box-shadow: 0 10px 24px rgba(124, 58, 237, 0.25);
}

.pagination .nav-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 14px 28px rgba(124, 58, 237, 0.35);
}

.pagination .nav-btn:disabled {
  background: linear-gradient(135deg, rgba(168,85,247,.5) 0%, rgba(124,58,237,.5) 100%);
  box-shadow: none;
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

/* 紧凑模式：与主页一致的侧边距策略 */
@media (max-width: 1330px) {
  .article-list-view { margin-top: 40px; }
  .header-section,
  .tabs-section,
  .content-section,
  .pagination-section {
    width: 66.666%;
    margin: 0 auto;
    padding-left: 0;
    padding-right: 0;
    min-width: 480px;
  }
  .header-section { margin-top: 0; }
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
  .article-list-view { margin-top: 30px; }
  /* 手机端也保持主体 2/3 宽度与最小宽度 480px（与主页一致） */
  .header-section,
  .tabs-section,
  .content-section,
  .pagination-section {
    width: 66.666%;
    margin: 0 auto;
    padding-left: 0;
    padding-right: 0;
    min-width: 480px;
  }
  .header-section { margin-top: 0; }

  .header-image {
    height: 250px;
  }

  .typing-text {
    font-size: 1.8rem;
    padding: 20px;
  }

  .article-grid {
    grid-template-columns: 1fr; /* 恢复博客/随笔在手机端单列展示 */
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
