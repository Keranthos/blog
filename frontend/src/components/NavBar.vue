<template>
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
        <router-link to="/" class="menu-item" active-class="active">
          <font-awesome-icon icon="house" class="menu-icon" />
          <span class="menu-text">首页</span>
        </router-link>
        <router-link to="/blog" class="menu-item" active-class="active">
          <font-awesome-icon icon="blog" class="menu-icon" />
          <span class="menu-text">博客</span>
        </router-link>
        <router-link to="/moments" class="menu-item" active-class="active">
          <font-awesome-icon icon="pen-to-square" class="menu-icon" />
          <span class="menu-text">随笔</span>
        </router-link>

        <!-- 媒体下拉菜单 -->
        <div class="menu-item dropdown" :class="{ active: isMediaActive }" @mouseenter="showDropdown" @mouseleave="hideDropdown">
          <font-awesome-icon icon="ellipsis" class="menu-icon" />
          <span class="menu-text">媒体</span>
          <!-- 箭头 - 独立于过渡动画，立即显示，指向"媒体"文本 -->
          <div v-if="dropdownVisible" class="dropdown-menu-arrow media-arrow"></div>
          <transition name="dropdown-fade">
            <div v-if="dropdownVisible" class="dropdown-menu">
              <router-link to="/fragments/books" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="bars" />
                <span>书单</span>
              </router-link>
              <router-link to="/fragments/novels" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="bookmark" />
                <span>小说</span>
              </router-link>
              <router-link to="/fragments/movies" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="film" />
                <span>电影</span>
              </router-link>
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
              <router-link to="/questionbox" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="question" />
                <span>提问箱</span>
              </router-link>
              <router-link to="/timeline" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="clock" />
                <span>时间树</span>
              </router-link>
              <router-link to="/presentation" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="chalkboard" />
                <span>讲演</span>
              </router-link>
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
              <router-link to="/images" class="dropdown-item" active-class="active">
                <font-awesome-icon icon="images" />
                <span>图片管理</span>
              </router-link>
              <router-link to="/location-update" class="dropdown-item" active-class="active" @click="settingsDropdownVisible = false">
                <font-awesome-icon icon="location-dot" />
                <span>更新位置</span>
              </router-link>
            </div>
          </transition>
        </div>
      </div>

      <!-- 右侧功能区 -->
      <div class="navbar-actions">
        <!-- 搜索框 -->
        <div class="search-box" :class="{ expanded: searchExpanded }">
          <font-awesome-icon icon="magnifying-glass" class="search-icon" @mousedown.prevent="handleSearch" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="想找点什么？"
            class="search-input"
            @focus="searchExpanded = true"
            @blur="searchExpanded = false"
            @keydown="onSearchKeydown"
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
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useStore } from 'vuex'
import { useRouter, useRoute } from 'vue-router'
import { getAllComments } from '@/api/Comments/browse'

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

// 搜索功能
const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    // 如果搜索框为空，展开搜索框并聚焦
    searchExpanded.value = true
    return
  }

  // 跳转到搜索结果页面
  router.push({
    path: '/search',
    query: { search: searchQuery.value }
  })
  searchQuery.value = ''
  searchExpanded.value = false
}

// 监听搜索框回车
const onSearchKeydown = (e) => {
  if (e.key === 'Enter') {
    handleSearch()
  }
}

// 监听路由变化，清理所有下拉框状态
watch(() => route.path, () => {
  dropdownVisible.value = false
  otherDropdownVisible.value = false
  settingsDropdownVisible.value = false
  showCreateMenu.value = false
  showCommentsDropdown.value = false
  searchExpanded.value = false

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
  border: 2px solid #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.brand-avatar:hover {
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.brand-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.brand-name {
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  background: rgba(102, 126, 234, 0.08);
  color: #667eea;
  animation: bounce 0.6s ease;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
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
  background: #667eea;
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  transform: translateY(-1px);
}

.menu-item.active:hover {
  background: #5a67d8;
  animation: bounce 0.6s ease;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
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
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 6px;
  padding: 4px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  z-index: 2000;
}

/* 下拉菜单箭头 - 独立于过渡动画，立即显示，指向菜单项文本 */
.dropdown-menu-arrow {
  position: absolute;
  top: calc(100% + 10px);
  width: 12px;
  height: 12px;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
  box-shadow: -2px -2px 4px rgba(0, 0, 0, 0.1);
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
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
}

.dropdown-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
  color: #667eea;
  font-weight: 600;
  border-left: 3px solid #667eea;
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
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  z-index: 2000;
}

/* 箭头 - 独立于过渡动画，立即显示，指向"写点什么"文本 */
.create-menu-arrow {
  position: absolute;
  top: calc(100% + 10px); /* 下拉框在100%+16px，箭头在下拉框上方6px，所以是16-6=10px */
  right: calc(65% - 15px); /* 菜单项文本中心大约在35%位置，从右边计算是65%，再向右移动20px */
  width: 12px;
  height: 12px;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
  box-shadow: -2px -2px 4px rgba(0, 0, 0, 0.1);
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
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
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
  background: rgba(102, 126, 234, 0.08);
  border-radius: 12px;
  padding: 10px 16px;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.search-box.expanded {
  background: white;
  border-color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
}

.search-icon {
  color: #667eea;
  font-size: 1rem;
  transition: all 0.3s ease;
  cursor: pointer;
}

.search-icon:hover {
  color: #764ba2;
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

.search-box.expanded .search-input {
  width: 200px;
  margin-left: 10px;
}

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
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-bottom: none;
  border-right: none;
  transform: rotate(45deg);
  transform-origin: center center;
  z-index: 2001;
  pointer-events: none;
  box-shadow: -2px -2px 4px rgba(0, 0, 0, 0.1);
  /* 箭头相对于评论菜单项定位，在下拉框顶部上方6px */
}

/* 下拉框包装器 - 包含箭头和下拉框 */
.comments-menu-wrapper {
  position: absolute;
  top: calc(100% + 16px);
  right: -20px; /* 评论下拉框向右移动20px */
  z-index: 2000;
  transform-origin: top center;
}

.comments-menu {
  width: 500px;
  max-height: 500px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.3);
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
  border: 3px solid rgba(102, 126, 234, 0.1);
  border-top-color: #667eea;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  color: #667eea;
  flex-shrink: 0;
}

.article-title-text {
  font-size: 1rem; /* 增大字体 */
  font-weight: 500;
  color: #667eea;
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
  background: rgba(102, 126, 234, 0.3) !important;
  border-radius: 3px !important;
  display: block !important;
}

.comments-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(102, 126, 234, 0.5) !important;
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
  border: 2px solid #667eea;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.user-avatar:hover {
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-prompt {
  background: rgba(102, 126, 234, 0.1);
  border-color: #667eea;
}

.login-prompt span {
  color: #667eea;
  font-size: 0.85rem;
  font-weight: 600;
}

.login-prompt:hover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
