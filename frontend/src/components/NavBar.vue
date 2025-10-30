<template>
  <div>
    <nav class="navbar">
      <div class="navbar-container">
        <!-- 左侧 Logo 区域 -->
        <div class="navbar-brand" @click="goToProfile">
          <div class="brand-avatar">
            <img :src="require('@/assets/my_headportrait.jpg')" alt="Avatar" />
          </div>
          <span class="brand-name">山角函兽</span>
        </div>

        <!-- 中间导航菜单 -->
        <div class="navbar-menu-center">
          <div class="menu-item" :class="{ active: route.path === '/' }" @click="navigateToHome">
            <font-awesome-icon icon="house" class="menu-icon" />
            <span class="menu-text">首页</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/blog' }" @click="navigateToBlog">
            <font-awesome-icon icon="blog" class="menu-icon" />
            <span class="menu-text">博客</span>
          </div>
          <div class="menu-item" :class="{ active: route.path === '/moments' }" @click="navigateToMoments">
            <font-awesome-icon icon="pen-to-square" class="menu-icon" />
            <span class="menu-text">随笔</span>
          </div>

          <!-- 媒体下拉菜单 -->
          <div class="menu-item dropdown" :class="{ active: isMediaActive }" @mouseenter="showDropdown" @mouseleave="hideDropdown">
            <font-awesome-icon icon="ellipsis" class="menu-icon" />
            <span class="menu-text">媒体</span>
            <!-- 箭头 - 独立于过渡动画，立即显示，指向"媒体"文本 -->
            <div v-if="dropdownVisible" class="dropdown-menu-arrow media-arrow"></div>
            <transition name="dropdown-fade">
              <div v-if="dropdownVisible" class="dropdown-menu">
                <div class="dropdown-item" :class="{ active: route.path === '/fragments/books' }" @click="navigateToBooks">
                  <font-awesome-icon icon="bars" />
                  <span>书单</span>
                </div>
                <div class="dropdown-item" :class="{ active: route.path === '/fragments/novels' }" @click="navigateToNovels">
                  <font-awesome-icon icon="bookmark" />
                  <span>小说</span>
                </div>
                <div class="dropdown-item" :class="{ active: route.path === '/fragments/movies' }" @click="navigateToMovies">
                  <font-awesome-icon icon="film" />
                  <span>电影</span>
                </div>
              </div>
            </transition>
          </div>

          <!-- 其他下拉菜单 -->
          <div class="menu-item dropdown" :class="{ active: isOtherActive }" @mouseenter="showOtherDropdown" @mouseleave="hideOtherDropdown">
            <font-awesome-icon icon="bars" class="menu-icon" />
            <span class="menu-text">其他</span>
            <!-- 箭头 - 独立于过渡动画，立即显示，指向"其他"文本 -->
            <div v-if="otherDropdownVisible" class="dropdown-menu-arrow other-arrow"></div>
            <transition name="dropdown-fade">
              <div v-if="otherDropdownVisible" class="dropdown-menu">
                <div class="dropdown-item" :class="{ active: route.path === '/questionbox' }" @click="navigateToQuestionbox">
                  <font-awesome-icon icon="question" />
                  <span>提问箱</span>
                </div>
                <div class="dropdown-item" :class="{ active: route.path === '/timeline' }" @click="navigateToTimeline">
                  <font-awesome-icon icon="clock" />
                  <span>时间树</span>
                </div>
                <div class="dropdown-item" :class="{ active: route.path === '/presentation' }" @click="navigateToPresentation">
                  <font-awesome-icon icon="chalkboard" />
                  <span>讲演</span>
                </div>
              </div>
            </transition>
          </div>

          <!-- 设置菜单（仅管理员可见） -->
          <div v-if="userLevel >= 3" class="menu-item dropdown" :class="{ active: isSettingsActive }" @mouseenter="showSettingsDropdown" @mouseleave="hideSettingsDropdown">
            <font-awesome-icon icon="gear" class="menu-icon" />
            <span class="menu-text">设置</span>
            <!-- 箭头 - 独立于过渡动画，立即显示，指向"设置"文本 -->
            <div v-if="settingsDropdownVisible" class="dropdown-menu-arrow settings-arrow"></div>
            <transition name="dropdown-fade">
              <div v-if="settingsDropdownVisible" class="dropdown-menu">
                <div class="dropdown-item" :class="{ active: route.path === '/images' }" @click="navigateToImages">
                  <font-awesome-icon icon="images" />
                  <span>图片管理</span>
                </div>
                <div class="dropdown-item" :class="{ active: route.path === '/location-update' }" @click="navigateToLocationUpdate">
                  <font-awesome-icon icon="location-dot" />
                  <span>更新位置</span>
                </div>
              </div>
            </transition>
          </div>
        </div>

        <!-- 右侧功能区 -->
        <div class="navbar-actions">
          <!-- 搜索键（按钮样式） -->
          <div class="search-box" :class="{ expanded: searchExpanded }" @click="openSearchModal">
            <font-awesome-icon icon="magnifying-glass" class="search-icon" />
            <span class="search-placeholder">找点什么？</span>
            <input
              v-model="searchQuery"
              type="text"
              class="search-input search-input-hidden"
              tabindex="-1"
              aria-hidden="true"
            />
          </div>

          <!-- 创建内容下拉菜单 -->
          <div v-if="userLevel >= 3" class="create-dropdown menu-item dropdown" @mouseenter="showCreateMenuHandler" @mouseleave="hideCreateMenu">
            <svg class="menu-icon write-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 20h9" />
              <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
            </svg>
            <span class="menu-text">写点什么</span>
            <!-- 箭头 - 独立于过渡动画，立即显示，指向"写点什么"文本 -->
            <div v-if="showCreateMenu" class="create-menu-arrow"></div>
            <transition name="dropdown-fade">
              <div v-if="showCreateMenu" class="create-menu">
                <div class="menu-section">
                  <div class="menu-section-title">📝 文章</div>
                  <button class="create-item" @click="createContent('article', 'blog')">
                    <font-awesome-icon icon="blog" />
                    <span>博客</span>
                  </button>
                  <button class="create-item" @click="createContent('article', 'moment')">
                    <font-awesome-icon icon="comment-dots" />
                    <span>随笔</span>
                  </button>
                </div>
                <div class="menu-divider"></div>
                <div class="menu-section">
                  <div class="menu-section-title">🎬 媒体卡片</div>
                  <button class="create-item" @click="createContent('media', 'books')">
                    <font-awesome-icon icon="bars" />
                    <span>书单</span>
                  </button>
                  <button class="create-item" @click="createContent('media', 'novels')">
                    <font-awesome-icon icon="bookmark" />
                    <span>小说</span>
                  </button>
                  <button class="create-item" @click="createContent('media', 'movies')">
                    <font-awesome-icon icon="film" />
                    <span>电影</span>
                  </button>
                </div>
                <div class="menu-divider"></div>
                <div class="menu-section">
                  <div class="menu-section-title">📊 其他</div>
                  <button class="create-item" @click="createContent('presentation', 'ppt')">
                    <font-awesome-icon icon="chalkboard" />
                    <span>讲演</span>
                  </button>
                </div>
              </div>
            </transition>
          </div>

          <!-- 评论下拉菜单 -->
          <div class="comments-dropdown menu-item dropdown" @mouseenter="showCommentsMenu" @mouseleave="hideCommentsMenu">
            <svg class="menu-icon comment-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
              <line x1="8" y1="10" x2="8" y2="10.01" />
              <line x1="12" y1="10" x2="12" y2="10.01" />
              <line x1="16" y1="10" x2="16" y2="10.01" />
            </svg>
            <span class="menu-text">评论</span>
            <!-- 箭头 - 独立于过渡动画，立即显示 -->
            <div v-if="showCommentsDropdown" class="comments-menu-arrow"></div>
            <transition name="dropdown-fade">
              <div v-if="showCommentsDropdown" class="comments-menu-wrapper" @mouseenter="showCommentsMenu" @mouseleave="hideCommentsMenu">
                <div class="comments-menu">
                  <div v-if="commentsLoading" class="comments-loading">
                    <div class="spinner"></div>
                    <span>加载中...</span>
                  </div>
                  <div v-else-if="recentComments.length === 0" class="comments-empty">
                    <span>暂无评论</span>
                  </div>
                  <div v-else class="comments-list">
                    <div
                      v-for="comment in recentComments"
                      :key="comment.ID"
                      class="comment-item"
                      @click="goToCommentArticle(comment)"
                    >
                      <div class="comment-bubble">
                        <div class="comment-avatar-col">
                          <div class="avatar-square">{{ (comment.username || 'U').charAt(0).toUpperCase() }}</div>
                        </div>
                        <div class="comment-content-col">
                          <div class="comment-meta-row">
                            <span class="comment-author">{{ comment.username || '匿名用户' }}</span>
                            <span class="comment-time">{{ formatCommentTime(comment.CreatedAt) }}</span>
                          </div>
                          <div class="comment-article-row">
                            <font-awesome-icon :icon="getArticleTypeIcon(comment.articleType)" class="article-type-icon" />
                            <span class="article-title-text">
                              {{ comment.articleTitle || '未知文章' }}
                            </span>
                          </div>
                          <div class="comment-text">{{ stripMarkdown(comment.content || '') }}</div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
          </div>

          <!-- 用户头像 -->
          <div class="user-menu" @click="goToAuthUser">
            <div v-if="!user.isLogged" class="user-avatar login-prompt">
              <span>未登录</span>
            </div>
            <div v-else class="user-avatar">
              <img :src="user.avatar || defaultAvatar" alt="User Avatar" />
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- 搜索模态框 - 使用 Teleport 渲染到 body -->
    <Teleport to="body">
      <transition name="search-modal-fade">
        <div v-if="showSearchModal" class="search-modal-overlay" @click="closeSearchModal">
          <div class="search-modal-container" @click.stop>
            <div class="search-modal-input-wrapper">
              <input
                ref="searchModalInput"
                v-model="searchModalQuery"
                type="text"
                placeholder="请输入关键词进行搜索,空格隔开多个关键词"
                class="search-modal-input"
                @keydown.esc="closeSearchModal"
              />
            </div>

            <!-- 搜索结果区域 -->
            <div v-if="searchModalLoading || hasSearched || searchModalResults.length > 0" class="search-modal-results">
              <!-- 加载中 -->
              <div v-if="searchModalLoading" class="search-loading">
                <div class="loading-spinner"></div>
                <p>搜索中...</p>
              </div>

              <!-- 搜索结果列表 -->
              <div v-else-if="searchModalResults.length > 0" class="search-results-list">
                <div
                  v-for="result in searchModalResults"
                  :key="`${result.type}-${result.id}`"
                  class="search-result-item"
                  @click="goToSearchResult(result)"
                >
                  <div class="result-image">
                    <img :src="result.image || 'https://picsum.photos/id/1/400/300'" alt="封面" />
                  </div>
                  <div class="result-content">
                    <h3 class="result-title">
                      <font-awesome-icon :icon="getResultTypeIcon(result.type)" class="result-type-icon" />
                      {{ result.title }}
                    </h3>
                    <div class="result-tags-wrapper">
                      <div v-if="filteredTags(result.tags).length" class="result-tags">
                        <span v-for="tag in filteredTags(result.tags)" :key="tag" class="tag">{{ tag }}</span>
                      </div>
                    </div>
                    <div class="result-meta">
                      <span class="meta-item">
                        <font-awesome-icon icon="calendar" />
                        {{ formatResultDate(result.time) }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 无结果 -->
              <div v-else-if="hasSearched && !searchModalLoading && searchModalResults.length === 0" class="search-empty">
                <font-awesome-icon icon="magnifying-glass" class="empty-icon" />
                <p>没有找到相关内容</p>
              </div>
            </div>

            <!-- 初始提示（无搜索结果时显示） -->
            <div v-else class="search-modal-hint">
              <div class="search-modal-icon">
                <font-awesome-icon icon="magnifying-glass" />
              </div>
              <p>输入关键词以搜索</p>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, ref, watch, nextTick } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import { getAllComments } from '@/api/Comments/browse'
import { getArticlesList } from '@/api/Articles/browse'
import { getMomentsList } from '@/api/Moments/browse'

const store = useStore()
const router = useRouter()
const route = useRoute()
const user = computed(() => store.state.user)

const dropdownVisible = ref(false)
const otherDropdownVisible = ref(false)
const settingsDropdownVisible = ref(false)
const showCreateMenu = ref(false)
const showCommentsDropdown = ref(false)
const searchExpanded = ref(false)
const searchQuery = ref('')
const showSearchModal = ref(false)
const searchModalQuery = ref('')
const searchModalInput = ref(null)
const searchModalResults = ref([])
const searchModalLoading = ref(false)
const searchTimeout = ref(null)
const hasSearched = ref(false)
const recentComments = ref([])
const commentsLoading = ref(false)
let timeout = null
let otherTimeout = null
let settingsTimeout = null
let commentsTimeout = null

const generateDefaultAvatar = (username) => {
  const canvas = document.createElement('canvas')
  canvas.width = 40
  canvas.height = 40
  const context = canvas.getContext('2d')
  context.fillStyle = '#6a1b9a' // 紫色背景
  context.fillRect(0, 0, canvas.width, canvas.height)
  if (username) {
    context.fillStyle = '#fff'
    context.font = '20px Arial'
    context.textAlign = 'center'
    context.textBaseline = 'middle'
    context.fillText(username.charAt(0).toUpperCase(), canvas.width / 2, canvas.height / 2)
  } else {
    context.fillStyle = '#fff'
    context.font = '14px Arial'
    context.textAlign = 'center'
    context.textBaseline = 'middle'
    context.fillText('未登录', canvas.width / 2, canvas.height / 2)
  }
  return canvas.toDataURL('image/png')
}

const defaultAvatar = generateDefaultAvatar(user.value.name)

// 创建内容（文章或媒体卡片）
const createContent = (contentType, subType) => {
  if (!user.value.isLogged) {
    router.push('/login-register')
    return
  }

  if (contentType === 'article') {
    // 创建文章，通过 query 参数指定文章类型
    router.push({
      path: '/editarticle',
      query: {
        contentType: 'article',
        articleType: subType
      }
    })
  } else if (contentType === 'media') {
    // 创建媒体卡片，通过 query 参数指定媒体类型
    router.push({
      path: '/editarticle',
      query: {
        contentType: 'media',
        mediaType: subType
      }
    })
  } else if (contentType === 'presentation') {
    // 创建讲演PPT，跳转到讲演编辑页面
    router.push({
      path: '/editpresentation',
      query: {
        contentType: 'presentation',
        presentationType: subType
      }
    })
  }

  showCreateMenu.value = false
}

const goToAuthUser = () => {
  if (user.value.isLogged) {
    router.push('/user-info')
  } else {
    router.push('/login-register')
  }
}

const goToProfile = () => {
  router.push('/profile')
}

// 导航函数
const navigateToHome = () => {
  router.push('/')
}

const navigateToBlog = () => {
  router.push('/blog')
}

const navigateToMoments = () => {
  router.push('/moments')
}

const navigateToBooks = () => {
  router.push('/fragments/books')
}

const navigateToNovels = () => {
  router.push('/fragments/novels')
}

const navigateToMovies = () => {
  router.push('/fragments/movies')
}

const navigateToQuestionbox = () => {
  router.push('/questionbox')
}

const navigateToTimeline = () => {
  router.push('/timeline')
}

const navigateToPresentation = () => {
  router.push('/presentation')
}

const navigateToImages = () => {
  router.push('/images')
}

const navigateToLocationUpdate = () => {
  router.push('/location-update')
  settingsDropdownVisible.value = false
}

const userLevel = computed(() => {
  if (user.value.isLogged) {
    return user.value.level
  } else {
    return 0
  }
})

// 检查当前路由是否匹配下拉框中的路径
const isMediaActive = computed(() => {
  const mediaPaths = ['/fragments/books', '/fragments/novels', '/fragments/movies']
  return mediaPaths.includes(route.path)
})

const isOtherActive = computed(() => {
  const otherPaths = ['/questionbox', '/timeline', '/presentation']
  return otherPaths.includes(route.path)
})

const isSettingsActive = computed(() => {
  const settingsPaths = ['/images', '/location-update']
  return settingsPaths.includes(route.path)
})

const showDropdown = () => {
  clearTimeout(timeout)
  // 关闭其他下拉框
  otherDropdownVisible.value = false
  settingsDropdownVisible.value = false
  showCreateMenu.value = false
  showCommentsDropdown.value = false
  dropdownVisible.value = true
}

const hideDropdown = () => {
  timeout = setTimeout(() => {
    dropdownVisible.value = false
  }, 200)
}

const showOtherDropdown = () => {
  clearTimeout(otherTimeout)
  // 关闭其他下拉框
  dropdownVisible.value = false
  showCreateMenu.value = false
  settingsDropdownVisible.value = false
  showCommentsDropdown.value = false
  otherDropdownVisible.value = true
}

const hideOtherDropdown = () => {
  otherTimeout = setTimeout(() => {
    otherDropdownVisible.value = false
  }, 200)
}

const showSettingsDropdown = () => {
  clearTimeout(settingsTimeout)
  // 关闭其他下拉框
  dropdownVisible.value = false
  showCreateMenu.value = false
  otherDropdownVisible.value = false
  showCommentsDropdown.value = false
  settingsDropdownVisible.value = true
}

const hideSettingsDropdown = () => {
  settingsTimeout = setTimeout(() => {
    settingsDropdownVisible.value = false
  }, 200)
}

const showCreateMenuHandler = () => {
  // 关闭其他下拉框
  dropdownVisible.value = false
  otherDropdownVisible.value = false
  settingsDropdownVisible.value = false
  showCommentsDropdown.value = false
  showCreateMenu.value = true
}

const hideCreateMenu = () => {
  showCreateMenu.value = false
}

// 加载最近评论
const loadRecentComments = async () => {
  // 如果正在加载，不重复加载
  if (commentsLoading.value) return

  try {
    commentsLoading.value = true
    const res = await getAllComments()
    console.log('获取到的评论响应:', res)
    // 后端返回格式为 {data: [...]}，getAllComments 已经返回了 res.data，所以需要再取 .data
    const data = res?.data || []
    console.log('获取到的评论数据:', data)
    // 只取最近的评论，最多显示10条
    recentComments.value = (Array.isArray(data) ? data : []).slice(0, 10)
    console.log('处理后的评论数据:', recentComments.value)
  } catch (error) {
    console.error('加载评论失败:', error)
    recentComments.value = []
  } finally {
    // 确保加载状态被重置
    commentsLoading.value = false
  }
}

// 显示评论菜单
const showCommentsMenu = () => {
  clearTimeout(commentsTimeout)
  // 关闭其他下拉框
  dropdownVisible.value = false
  otherDropdownVisible.value = false
  settingsDropdownVisible.value = false
  showCreateMenu.value = false

  // 先显示下拉框
  showCommentsDropdown.value = true

  // 只有在没有评论数据且不在加载中时才加载
  if (recentComments.value.length === 0 && !commentsLoading.value) {
    loadRecentComments()
  }
}

// 隐藏评论菜单
const hideCommentsMenu = () => {
  clearTimeout(commentsTimeout)
  commentsTimeout = setTimeout(() => {
    showCommentsDropdown.value = false
    // 不清空评论数据，保留以便下次快速显示
  }, 200)
}

// 格式化评论时间
const formatCommentTime = (timestamp) => {
  if (!timestamp) return '未知时间'
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 去除Markdown语法
const stripMarkdown = (text) => {
  if (!text) return ''
  // 移除Markdown语法标记
  return text
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1') // 链接
    .replace(/!\[([^\]]*)\]\([^)]*\)/g, '$1') // 图片
    .replace(/#{1,6}\s+/g, '') // 标题
    .replace(/\*\*([^*]+)\*\*/g, '$1') // 粗体
    .replace(/\*([^*]+)\*/g, '$1') // 斜体
    .replace(/`([^`]+)`/g, '$1') // 行内代码
    .replace(/```[\s\S]*?```/g, '') // 代码块
    .replace(/>\s+/g, '') // 引用
    .replace(/^\s*[-*+]\s+/gm, '') // 列表
    .replace(/^\s*\d+\.\s+/gm, '') // 有序列表
    .trim()
    .substring(0, 100) // 限制长度
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

// 跳转到评论所属的文章
const goToCommentArticle = (comment) => {
  if (!comment.articleType || !comment.blogID) return

  let path = '/'
  if (comment.articleType === 'moment') {
    path = `/moments/${comment.blogID}`
  } else if (comment.articleType === 'blog') {
    path = `/blog/${comment.blogID}`
  } else if (comment.articleType === 'research') {
    path = `/research/${comment.blogID}`
  } else if (comment.articleType === 'project') {
    path = `/project/${comment.blogID}`
  }

  router.push(path)
  showCommentsDropdown.value = false
}

// 点击按钮直接打开模态搜索

// 已改为点击按钮打开模态搜索，不再监听输入框回车

// 打开搜索模态框
const openSearchModal = () => {
  showSearchModal.value = true
  searchModalQuery.value = ''
  nextTick(() => {
    if (searchModalInput.value) {
      searchModalInput.value.focus()
    }
  })
}

// 关闭搜索模态框
const closeSearchModal = () => {
  showSearchModal.value = false
  searchModalQuery.value = ''
  searchModalResults.value = []
  hasSearched.value = false
  // 清除搜索定时器
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
    searchTimeout.value = null
  }
}

// 获取结果类型图标
const getResultTypeIcon = (type) => {
  const iconMap = {
    blog: 'blog',
    moment: 'pen-to-square',
    research: 'flask',
    project: 'code'
  }
  return iconMap[type] || 'file'
}

// 过滤搜索结果中的标签，移除类型性标签（如 博客/随笔/研究/项目 以及英文类型）
const filteredTags = (tags = []) => {
  if (!Array.isArray(tags)) return []
  const typeNames = new Set(['博客', '随笔', '研究', '项目', 'blog', 'moment', 'research', 'project'])
  return tags.filter(t => {
    const s = String(t).trim()
    return s && !typeNames.has(s.toLowerCase ? s.toLowerCase() : s)
  })
}

// 格式化日期
const formatResultDate = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 跳转到搜索结果文章
const goToSearchResult = (result) => {
  if (result.type === 'moment') {
    router.push(`/moments/${result.id}`)
  } else if (result.type === 'blog') {
    router.push(`/blog/${result.id}`)
  } else if (result.type === 'research') {
    router.push(`/research/${result.id}`)
  } else if (result.type === 'project') {
    router.push(`/project/${result.id}`)
  }
  closeSearchModal()
}

// 检查文章是否匹配关键词（支持空格分隔的多个关键词）
const matchesKeywords = (item, keywords) => {
  const title = (item.title || '').toLowerCase()
  const tags = (item.tags || []).join(' ').toLowerCase()

  // 只要任意一个关键词匹配标题或标签即可
  return keywords.some(keyword =>
    title.includes(keyword.toLowerCase()) ||
    tags.includes(keyword.toLowerCase())
  )
}

// 模态框搜索
const handleModalSearch = async () => {
  const query = searchModalQuery.value.trim()
  if (!query) {
    searchModalResults.value = []
    hasSearched.value = false
    return
  }

  // 分割关键词（支持空格分隔）
  const keywords = query.split(/\s+/).filter(k => k.length > 0)
  if (keywords.length === 0) {
    searchModalResults.value = []
    hasSearched.value = false
    return
  }

  searchModalLoading.value = true
  searchModalResults.value = []
  hasSearched.value = true

  try {
    const allResults = []
    const types = ['blog', 'moment', 'research', 'project']

    for (const type of types) {
      try {
        if (type === 'moment') {
          const response = await getMomentsList(1, 100)
          const articles = (response.data || []).filter(item =>
            matchesKeywords(item, keywords)
          ).map(item => ({
            id: item.ID,
            type: 'moment',
            title: item.title,
            image: item.image,
            tags: item.tags || [],
            time: item.CreatedAt
          }))
          allResults.push(...articles)
        } else {
          const response = await getArticlesList(type, 1, 100)
          const articles = (response.data || []).filter(item =>
            matchesKeywords(item, keywords)
          ).map(item => ({
            id: item.ID,
            type,
            title: item.title,
            image: item.image,
            tags: item.tags || [],
            time: item.CreatedAt
          }))
          allResults.push(...articles)
        }
      } catch (error) {
        console.error(`搜索 ${type} 失败:`, error)
      }
    }

    searchModalResults.value = allResults
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    searchModalLoading.value = false
  }
}

// 监听搜索输入，1秒后自动搜索
watch(searchModalQuery, (newQuery) => {
  // 清除之前的定时器
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
    searchTimeout.value = null
  }

  // 如果输入为空，清空结果和搜索状态
  if (!newQuery.trim()) {
    searchModalResults.value = []
    hasSearched.value = false
    return
  }

  // 设置1秒延迟搜索
  searchTimeout.value = setTimeout(() => {
    handleModalSearch()
  }, 1000)
})

// 监听路由变化，清理所有下拉框状态
watch(() => route.path, () => {
  dropdownVisible.value = false
  otherDropdownVisible.value = false
  settingsDropdownVisible.value = false
  showCreateMenu.value = false
  showCommentsDropdown.value = false
  searchExpanded.value = false
  showSearchModal.value = false
  searchModalResults.value = []
  hasSearched.value = false

  // 清理所有timeout
  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  }
  if (otherTimeout) {
    clearTimeout(otherTimeout)
    otherTimeout = null
  }
  if (settingsTimeout) {
    clearTimeout(settingsTimeout)
    settingsTimeout = null
  }
  if (commentsTimeout) {
    clearTimeout(commentsTimeout)
    commentsTimeout = null
  }
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
    searchTimeout.value = null
  }
})

</script>

<style scoped>
/* 导航栏容器 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
  padding: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.7) 0%, rgba(248, 250, 252, 0.6) 100%);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.navbar-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0px 32px;
}

/* 左侧品牌区域 */
.navbar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: opacity 0.3s ease;
  flex-shrink: 0;
  margin-left: -16px;
  padding-left: 16px;
}

.navbar-brand:hover {
  opacity: 0.8;
}

.brand-avatar {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #a855f7;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.3);
  transition: all 0.3s ease;
}

.brand-avatar:hover {
  box-shadow: 0 6px 20px rgba(168, 85, 247, 0.5);
}

.brand-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.brand-name {
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 0.5px;
}

/* 中间菜单区域 */
.navbar-menu-center {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  justify-content: flex-start;
  margin-left: 40px;
}

.menu-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  color: #64748b;
  text-decoration: none;
  border-radius: 12px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.menu-icon {
  font-size: 1.2rem;
  transition: all 0.3s ease;
}

/* 简洁图标样式 */
.write-icon,
.comment-icon {
  width: 1.2rem;
  height: 1.2rem;
  stroke: currentColor;
}

.menu-text {
  font-size: 0.85rem;
  transition: all 0.3s ease;
}

/* Hover 效果 */
.menu-item:hover {
  background: rgba(168, 85, 247, 0.08);
  color: #a855f7;
  animation: bounce 0.6s ease;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.2);
}

/* 跳一跳动画 */
@keyframes bounce {
  0% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
  100% {
    transform: translateY(0);
  }
}

/* 激活状态 */
.menu-item.active {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  box-shadow: 0 2px 8px rgba(168, 85, 247, 0.2);
  transform: translateY(-1px);
}

.menu-item.active:hover {
  background: rgba(168, 85, 247, 0.15);
  animation: bounce 0.6s ease;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.3);
}

/* 下拉菜单 */
.dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 16px);
  left: 0;
  min-width: 160px;
  max-width: 200px;
  background: rgba(255, 255, 255, 1);
  border-radius: 6px;
  padding: 4px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-top: none;
  z-index: 2003;
}

/* 下拉菜单箭头 - 独立于过渡动画，立即显示，指向菜单项文本 */
.dropdown-menu-arrow {
  position: absolute;
  top: calc(100% + 10px);
  width: 12px;
  height: 12px;
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
}

/* 媒体箭头 - 指向"媒体"文本中心 */
.media-arrow {
  left: 35%;
}

/* 其他箭头 - 指向"其他"文本中心 */
.other-arrow {
  left: 35%;
}

/* 设置箭头 - 指向"设置"文本中心 */
.settings-arrow {
  left: 35%;
}

.dropdown-item {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 10px;
  padding: 10px 14px;
  color: #5a5a5a;
  text-decoration: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
  white-space: nowrap;
  text-align: left;
}

.dropdown-item:hover {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%);
  color: #a855f7;
}

.dropdown-item.active {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.15) 0%, rgba(124, 58, 237, 0.15) 100%);
  color: #a855f7;
  font-weight: 600;
  border-left: 3px solid #a855f7;
}

.dropdown-item svg {
  font-size: 0.9rem;
}

/* 下拉菜单过渡动画 - 从扁平压缩逐渐撑开到正常高度 */
.dropdown-fade-enter-active {
  transition: transform 0.35s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.25s ease;
  overflow: hidden;
}

.dropdown-fade-leave-active {
  transition: transform 0.2s cubic-bezier(0.4, 0, 1, 1), opacity 0.15s ease;
  overflow: hidden;
}

.dropdown-fade-enter-from {
  opacity: 1;
  transform: scaleY(0);
  transform-origin: top center;
  /* 初始状态：所有内容都被压缩成扁平（垂直方向缩放为0） */
}

.dropdown-fade-enter-to {
  opacity: 1;
  transform: scaleY(1);
  transform-origin: top center;
  /* 最终状态：完全展开到正常高度 */
}

.dropdown-fade-leave-from {
  opacity: 1;
  transform: scaleY(1);
  transform-origin: top center;
}

.dropdown-fade-leave-to {
  opacity: 0;
  transform: scaleY(0);
  transform-origin: top center;
}

/* 箭头在过渡动画中立即显示，不受 scaleY 影响 */
.dropdown-fade-enter-from .comments-menu-arrow,
.dropdown-fade-enter-active .comments-menu-arrow,
.dropdown-fade-enter-to .comments-menu-arrow,
.dropdown-fade-leave-from .comments-menu-arrow,
.dropdown-fade-leave-active .comments-menu-arrow,
.dropdown-fade-leave-to .comments-menu-arrow {
  opacity: 1 !important;
  transform: rotate(45deg) scaleY(1) !important;
  transform-origin: center center !important;
  transition: none !important;
  /* 确保箭头不受父元素的 scaleY 影响，强制 scaleY(1) */
}

/* 创建菜单箭头也立即显示 */
.dropdown-fade-enter-from .create-menu-arrow,
.dropdown-fade-enter-active .create-menu-arrow,
.dropdown-fade-enter-to .create-menu-arrow,
.dropdown-fade-leave-from .create-menu-arrow,
.dropdown-fade-leave-active .create-menu-arrow,
.dropdown-fade-leave-to .create-menu-arrow {
  opacity: 1 !important;
  transform: rotate(45deg) !important;
  transform-origin: center center !important;
  transition: none !important;
}

/* 创建内容下拉菜单 */
.create-dropdown {
  position: relative;
}

.create-menu {
  position: absolute;
  top: calc(100% + 16px);
  right: -20px; /* 右侧菜单项，下拉框向右移动20px */
  min-width: 200px;
  background: rgba(255, 255, 255, 1);
  border-radius: 8px;
  padding: 12px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-top: none;
  z-index: 2003;
}

/* 箭头 - 独立于过渡动画，立即显示，指向"写点什么"文本 */
.create-menu-arrow {
  position: absolute;
  top: calc(100% + 10px); /* 下拉框在100%+16px，箭头在下拉框上方6px，所以是16-6=10px */
  right: calc(65% - 15px); /* 菜单项文本中心大约在35%位置，从右边计算是65%，再向右移动20px */
  width: 12px;
  height: 12px;
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
  /* 箭头相对于"写点什么"菜单项定位，在下拉框顶部上方6px */
}

.menu-section {
  padding: 8px 0;
}

.menu-section-title {
  padding: 8px 12px;
  font-size: 0.85rem;
  font-weight: 700;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.menu-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, #e0e0e0, transparent);
  margin: 8px 0;
}

.create-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  color: #5a5a5a;
  background: transparent;
  border: none;
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all 0.2s ease;
  cursor: pointer;
  text-align: left;
}

.create-item:hover {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%);
  color: #a855f7;
}

.create-item svg {
  font-size: 1rem;
  width: 16px;
}

/* 右侧功能区 */
.navbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  margin-right: -16px;
  padding-right: 16px;
}

/* 搜索框 */
.search-box {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(168, 85, 247, 0.04);
  border-radius: 14px;
  padding: 8px 14px;
  transition: all 0.2s ease;
  border: 2px solid #a855f7;
  color: #a855f7;
  cursor: pointer;
}

.search-box:hover {
  background: rgba(168, 85, 247, 0.08);
}

.search-icon {
  color: #a855f7;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.search-placeholder {
  color: #a855f7;
  font-size: 0.95rem;
}

.search-input {
  border: none;
  outline: none;
  background: transparent;
  color: #333;
  font-size: 0.9rem;
  width: 0;
  margin-left: 0;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-input-hidden { display: none; }

.search-input::placeholder {
  color: #999;
}

/* 评论下拉菜单和创建内容下拉菜单现在使用 menu-item 样式 */
.comments-dropdown {
  position: relative;
}

/* 箭头 - 独立于过渡动画，立即显示，相对于评论菜单项定位 */
.comments-menu-arrow {
  position: absolute;
  top: calc(100% + 10px); /* 下拉框在100%+16px，箭头在下拉框上方6px，所以是16-6=10px */
  right: calc(65% - 15px); /* 菜单项文本中心大约在35%位置，从右边计算是65%，再向右移动20px */
  width: 12px;
  height: 12px;
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
  /* 箭头相对于评论菜单项定位，在下拉框顶部上方6px */
}

/* 下拉框包装器 - 包含箭头和下拉框 */
.comments-menu-wrapper {
  position: absolute;
  top: calc(100% + 16px);
  right: -20px; /* 评论下拉框向右移动20px */
  z-index: 2003;
  transform-origin: top center;
}

.comments-menu {
  width: 500px;
  max-height: 500px;
  background: rgba(255, 255, 255, 1);
  border-radius: 8px;
  padding: 16px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-top: none;
  overflow-y: auto !important; /* 强制启用滚动 */
  scrollbar-width: thin !important; /* Firefox */
}

.comments-loading,
.comments-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #999;
  gap: 12px;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(168, 85, 247, 0.1);
  border-top-color: #a855f7;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  cursor: pointer;
  transition: all 0.2s ease;
}

.comment-item:hover {
  transform: translateY(-2px);
}

.comment-item:hover .comment-bubble {
  background: #F0F0F0; /* 悬停时深灰色（原来的浅灰色） */
  transform: translateY(-1px);
}

.comment-bubble {
  background: #F8F8F8; /* 普通状态更浅的灰色 */
  border-radius: 5px;
  padding: 1.2rem 1.4rem; /* 增大内边距使卡片更大 */
  border-top: 1.5px solid rgba(200, 200, 200, 0.4);
  border-left: 1.5px solid rgba(200, 200, 200, 0.4);
  border-bottom: 1px solid rgba(200, 200, 200, 0.15);
  border-right: 1px solid rgba(200, 200, 200, 0.15);
  transition: all 0.2s ease;
  display: flex;
  gap: 16px; /* 增大间距 */
  align-items: flex-start;
}

.comment-avatar-col {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
}

.avatar-square {
  width: 48px; /* 增大头像尺寸 */
  height: 48px;
  border-radius: 8px;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px; /* 增大头像文字 */
  font-weight: 600;
}

.comment-content-col {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px; /* 增大内容间距 */
}

.comment-meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.comment-author {
  font-size: 1rem; /* 增大字体 */
  font-weight: 600;
  color: #2b2b2b;
}

.comment-time {
  font-size: 0.85rem; /* 增大字体 */
  color: #999;
}

.comment-article-row {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.article-type-icon {
  font-size: 1rem; /* 增大图标 */
  color: #a855f7;
  flex-shrink: 0;
}

.article-title-text {
  font-size: 1rem; /* 增大字体 */
  font-weight: 500;
  color: #a855f7;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: left;
}

.comment-text {
  font-size: 0.95rem; /* 增大字体 */
  color: #2b2b2b;
  line-height: 1.6;
  word-break: break-word;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  text-align: left;
}

/* 滚动条样式 - 使用 !important 覆盖全局隐藏 */
.comments-menu::-webkit-scrollbar {
  width: 6px !important;
  height: 6px !important;
  background: transparent !important;
  display: block !important; /* 覆盖全局的 display: none */
}

.comments-menu::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05) !important;
  border-radius: 3px !important;
  display: block !important;
}

.comments-menu::-webkit-scrollbar-thumb {
  background: rgba(168, 85, 247, 0.3) !important;
  border-radius: 3px !important;
  display: block !important;
}

.comments-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(168, 85, 247, 0.5) !important;
}

/* 确保滚动条可见 - 覆盖所有可能的全局样式 */
.comments-menu::-webkit-scrollbar-corner,
.comments-menu::-webkit-scrollbar-button,
.comments-menu::-webkit-scrollbar-track-piece {
  display: block !important;
}

/* 用户菜单 */
.user-menu {
  cursor: pointer;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #a855f7;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
}

.user-avatar:hover {
  box-shadow: 0 6px 20px rgba(168, 85, 247, 0.5);
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-prompt {
  background: rgba(168, 85, 247, 0.1);
  border-color: #a855f7;
}

.login-prompt span {
  color: #a855f7;
  font-size: 0.85rem;
  font-weight: 600;
}

.login-prompt:hover {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
}

.login-prompt:hover span {
  color: white;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .navbar-container {
    padding: 4px 24px;
  }

  .navbar-brand {
    margin-left: -12px;
    padding-left: 12px;
  }

  .navbar-actions {
    margin-right: -12px;
    padding-right: 12px;
  }

  .navbar-menu-center {
    gap: 6px;
    margin-left: 30px;
  }

  .menu-item {
    padding: 8px 14px;
    font-size: 0.9rem;
    gap: 4px;
  }

  .menu-text {
    font-size: 0.9rem;
  }

  .menu-icon {
    font-size: 0.95rem;
  }

  .brand-name {
    font-size: 1.1rem;
  }

  .navbar-actions {
    gap: 8px;
  }
}

/* 搜索模态框 */
.search-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 10000;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px) saturate(150%);
  -webkit-backdrop-filter: blur(30px) saturate(150%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

/* 搜索框弹出时，看板娘不虚化：将其层级提升到遮罩之上 */
:global(.waifu),
:global(.waifu-tips),
:global(#waifu) {
  z-index: 10001 !important;
}

.search-modal-container {
  width: 90%;
  max-width: 600px;
  max-height: 85vh;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  padding: 2rem;
  border: 1px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  cursor: default;
  margin: 7.5vh auto;
}

.search-modal-input-wrapper {
  width: 100%;
  margin-bottom: 1.5rem;
  flex-shrink: 0;
}

.search-modal-input {
  width: 100%;
  padding: 1rem 1.5rem;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  font-size: 1rem;
  outline: none;
  transition: all 0.3s ease;
  background: white;
}

.search-modal-input:focus {
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.1);
}

.search-modal-icon {
  font-size: 3rem;
  color: #a855f7;
  margin: 1rem 0;
}

.search-modal-hint {
  text-align: center;
  color: #666;
  line-height: 1.8;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.search-modal-hint p {
  margin: 0.5rem 0;
  font-size: 0.95rem;
}

.search-modal-hint p:first-child {
  font-size: 1rem;
  color: #333;
  font-weight: 500;
}

/* 搜索结果区域 */
.search-modal-results {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 0;
  max-height: calc(85vh - 150px);
  margin-top: 1rem;
}

/* 隐藏滚动条但保持滚动功能 */
.search-modal-results {
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}

.search-modal-results::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

/* 加载中状态 */
.search-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(168, 85, 247, 0.1);
  border-top-color: #a855f7;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.search-loading p {
  font-size: 0.95rem;
  margin: 0;
}

/* 搜索结果列表 */
.search-results-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.search-result-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-top: 1px solid rgba(0, 0, 0, 0.12);
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-result-item:hover {
  background: rgba(168, 85, 247, 0.08);
  border-color: #a855f7;
}

.result-image {
  width: 120px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.result-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.result-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  min-width: 0;
}

.result-type {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  width: fit-content;
}

.result-title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #333;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.result-type-icon {
  color: #a855f7;
  margin-right: 8px;
}

.result-tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.result-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.result-tags .tag {
  padding: 0.25rem 0.6rem;
  background: rgba(168, 85, 247, 0.08);
  color: #a855f7;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.result-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 0.5rem;
  font-size: 0.85rem;
  color: #999;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

/* 无结果状态 */
.search-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  color: #999;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.search-empty p {
  margin: 0;
  font-size: 0.95rem;
}

/* 搜索模态框过渡动画 */
.search-modal-fade-enter-active,
.search-modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.search-modal-fade-enter-active .search-modal-container,
.search-modal-fade-leave-active .search-modal-container {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-modal-fade-enter-from {
  opacity: 0;
}

.search-modal-fade-enter-from .search-modal-container {
  transform: scale(0.9) translateY(-20px);
  opacity: 0;
}

.search-modal-fade-leave-to {
  opacity: 0;
}

.search-modal-fade-leave-to .search-modal-container {
  transform: scale(0.9) translateY(-20px);
  opacity: 0;
}

@media (max-width: 768px) {
  .navbar-container {
    padding: 4px 16px;
  }

  .navbar-brand {
    margin-left: -12px;
    padding-left: 12px;
  }

  .navbar-actions {
    margin-right: -12px;
    padding-right: 12px;
  }

  .navbar-menu-center {
    gap: 4px;
    margin-left: 20px;
  }

  .menu-item {
    padding: 8px 10px;
    font-size: 0.85rem;
    gap: 4px;
  }

  .menu-text {
    display: none;
  }

  .menu-icon {
    font-size: 1.1rem;
  }

  .brand-name {
    font-size: 1rem;
  }

  .search-box {
    padding: 8px 12px;
  }

  .search-box.expanded .search-input {
    width: 120px;
  }

  .navbar-actions {
    gap: 6px;
  }

  .action-btn span {
    display: none;
  }

}

/* 滚动时导航栏效果 */
@media (min-width: 769px) {
  .navbar.scrolled .navbar-container {
    padding: 8px 20px;
    background: rgba(255, 255, 255, 0.98);
  }
}
</style>
