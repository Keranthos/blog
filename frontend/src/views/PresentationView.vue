<template>
  <div class="presentation-view">
    <NavBar />

    <!-- 英雄图片区域 -->
    <div class="hero-section">
      <div class="hero-image">
        <img :src="heroImage" alt="Presentation Hero" />
        <div class="hero-overlay">
          <div class="hero-content">
            <h1 class="hero-title">讲演展示</h1>
            <p class="hero-subtitle">分享知识与见解</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 分隔文字 -->
    <div class="separator-section">
      <div class="separator-line"></div>
      <p class="separator-text">表达是人的本能，要如何满足和抑制它呢？</p>
    </div>

    <!-- 内容卡片区域 -->
    <div class="content-section">
      <div class="content-container">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>

        <!-- 卡片网格 -->
        <div v-else class="cards-grid">
          <div
            v-for="presentation in currentPagePresentations"
            :key="presentation.id"
            class="presentation-card"
            @click="openPresentation(presentation)"
          >
            <div class="card-image">
              <img
                :src="presentation.thumbnail || defaultThumbnail"
                :alt="presentation.title"
                @error="handleImageError"
              />
            </div>
            <div class="card-content">
              <h3 class="card-title">{{ presentation.title }}</h3>
              <div class="card-tags">
                <span
                  v-for="tag in presentation.tags"
                  :key="tag"
                  class="tag"
                >
                  {{ tag }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页控件 -->
        <div v-if="totalPages > 1" class="pagination">
          <button
            class="pagination-btn"
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
          >
            <font-awesome-icon icon="chevron-left" />
          </button>

          <div class="pagination-pages">
            <button
              v-for="page in visiblePages"
              :key="page"
              class="page-btn"
              :class="{ active: page === currentPage }"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
          </div>

          <button
            class="pagination-btn"
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
          >
            <font-awesome-icon icon="chevron-right" />
          </button>
        </div>
      </div>
    </div>

    <!-- PPT预览模态框 -->
    <div v-if="showPreview" class="preview-modal" @click="closePreview">
      <div class="preview-content" @click.stop>
        <button class="close-btn" @click="closePreview">
          <font-awesome-icon icon="times" />
        </button>
        <div class="preview-container">
          <iframe
            v-if="currentPresentation && currentPresentation.url"
            :src="currentPresentation.url"
            class="preview-iframe"
          ></iframe>
          <div v-else class="preview-placeholder">
            <font-awesome-icon icon="file-powerpoint" />
            <p>PPT预览功能开发中...</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/NavBar.vue'

// 响应式数据
const loading = ref(true)
const presentations = ref([])
const currentPage = ref(1)
const itemsPerPage = 9 // 每页9张卡片（3行×3列）
const showPreview = ref(false)
const currentPresentation = ref(null)

// 英雄图片
const heroImage = ref('/src/assets/background-image.jpg')
const defaultThumbnail = ref('/src/assets/background-image.jpg')

// 计算属性
const totalPages = computed(() => Math.ceil(presentations.value.length / itemsPerPage))

const currentPagePresentations = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return presentations.value.slice(start, end)
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages
})

// 方法
const loadPresentations = async () => {
  try {
    loading.value = true
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 模拟数据
    presentations.value = [
      {
        id: 1,
        title: 'Vue.js 3.0 新特性详解',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['Vue.js', '前端', 'JavaScript']
      },
      {
        id: 2,
        title: '人工智能在Web开发中的应用',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['AI', 'Web开发', '机器学习']
      },
      {
        id: 3,
        title: '微服务架构设计与实践',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['微服务', '架构', '后端']
      },
      {
        id: 4,
        title: 'React Hooks 深度解析',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['React', 'Hooks', '前端']
      },
      {
        id: 5,
        title: '数据库优化策略',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['数据库', '优化', '性能']
      },
      {
        id: 6,
        title: 'Docker 容器化部署',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['Docker', '容器', '部署']
      },
      {
        id: 7,
        title: 'TypeScript 高级特性',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['TypeScript', '类型', '前端']
      },
      {
        id: 8,
        title: 'GraphQL 实战指南',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['GraphQL', 'API', '后端']
      },
      {
        id: 9,
        title: 'Node.js 性能优化',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['Node.js', '性能', '优化']
      },
      {
        id: 10,
        title: 'CSS Grid 布局详解',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['CSS', 'Grid', '布局']
      },
      {
        id: 11,
        title: 'Webpack 5 配置指南',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['Webpack', '构建', '前端']
      },
      {
        id: 12,
        title: 'Redis 缓存策略',
        thumbnail: '/src/assets/background-image.jpg',
        url: '#',
        tags: ['Redis', '缓存', '数据库']
      }
    ]
  } catch (error) {
    console.error('加载讲演数据失败:', error)
  } finally {
    loading.value = false
  }
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const openPresentation = (presentation) => {
  currentPresentation.value = presentation
  showPreview.value = true
}

const closePreview = () => {
  showPreview.value = false
  currentPresentation.value = null
}

const handleImageError = (event) => {
  event.target.src = defaultThumbnail.value
}

onMounted(() => {
  loadPresentations()
})
</script>

<style scoped>
.presentation-view {
  min-height: 100vh;
  position: relative;
}

/* 英雄图片区域 */
.hero-section {
  position: relative;
  height: 400px;
  overflow: hidden;
}

.hero-image {
  position: relative;
  width: 100%;
  height: 100%;
}

.hero-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.1) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-content {
  text-align: center;
  color: white;
}

.hero-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 10px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.hero-subtitle {
  font-size: 1.2rem;
  margin: 0;
  opacity: 0.9;
}

/* 分隔区域 */
.separator-section {
  padding: 40px 0;
  text-align: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.separator-line {
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  margin: 0 auto 20px;
  border-radius: 2px;
}

.separator-text {
  font-size: 1.1rem;
  color: #333;
  margin: 0;
  font-style: italic;
}

/* 内容区域 */
.content-section {
  padding: 60px 0;
  background: transparent;
}

.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 加载状态 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 卡片网格 */
.cards-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 30px;
  margin-bottom: 60px;
}

.presentation-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.presentation-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.2);
}

.card-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.presentation-card:hover .card-image img {
  transform: scale(1.05);
}

.card-content {
  padding: 20px;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 15px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 0.8rem;
  font-weight: 500;
}

/* 分页控件 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 40px;
}

.pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #666;
}

.pagination-btn:hover:not(:disabled) {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-pages {
  display: flex;
  gap: 5px;
}

.page-btn {
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #666;
  font-weight: 500;
}

.page-btn:hover {
  background: #f5f5f5;
}

.page-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

/* 预览模态框 */
.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(5px);
}

.preview-content {
  position: relative;
  width: 90%;
  height: 90%;
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.close-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  width: 40px;
  height: 40px;
  border: none;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.7);
}

.preview-container {
  width: 100%;
  height: 100%;
}

.preview-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
  font-size: 1.2rem;
}

.preview-placeholder svg {
  font-size: 4rem;
  margin-bottom: 20px;
  color: #667eea;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .cards-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }

  .hero-title {
    font-size: 2.5rem;
  }
}

@media (max-width: 768px) {
  .cards-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .hero-section {
    height: 300px;
  }

  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .content-container {
    padding: 0 15px;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 5px;
  }

  .page-btn, .pagination-btn {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: 1.8rem;
  }

  .separator-text {
    font-size: 1rem;
  }

  .card-title {
    font-size: 1rem;
  }

  .tag {
    font-size: 0.75rem;
    padding: 3px 10px;
  }
}
</style>
