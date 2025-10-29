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
        <div v-if="userLevel >= 3" class="create-dropdown" @mouseenter="showCreateMenuHandler" @mouseleave="hideCreateMenu">
          <button class="action-btn edit-btn" title="写点什么">
            <font-awesome-icon icon="pen-to-square" />
            <span>写点什么</span>
          </button>
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

const store = useStore()
const router = useRouter()
const route = useRoute()
const user = computed(() => store.state.user)

const dropdownVisible = ref(false)
const otherDropdownVisible = ref(false)
const settingsDropdownVisible = ref(false)
const showCreateMenu = ref(false)
const searchExpanded = ref(false)
const searchQuery = ref('')
let timeout = null
let otherTimeout = null
let settingsTimeout = null

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
  showCreateMenu.value = true
}

const hideCreateMenu = () => {
  showCreateMenu.value = false
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

/* 添加箭头 - 位置偏右不在正中 */
.dropdown-menu::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 15%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-bottom: 6px solid rgba(255, 255, 255, 0.95);
  filter: drop-shadow(0 -2px 4px rgba(0, 0, 0, 0.1));
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

/* 下拉菜单过渡动画 */
.dropdown-fade-enter-active {
  transition: all 1.0s cubic-bezier(0.34, 1.56, 0.64, 1);
  overflow: hidden;
}

.dropdown-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 1, 1);
  overflow: hidden;
}

.dropdown-fade-enter-from {
  opacity: 0;
  max-height: 2px;
  clip-path: inset(0 0 100% 0);
}

.dropdown-fade-enter-to {
  max-height: 500px;
  clip-path: inset(0);
  opacity: 1;
}

.dropdown-fade-leave-from {
  max-height: 500px;
  clip-path: inset(0);
  opacity: 1;
}

.dropdown-fade-leave-to {
  opacity: 0;
  max-height: 2px;
  clip-path: inset(0 0 100% 0);
}

/* 创建内容下拉菜单 */
.create-dropdown {
  position: relative;
}

.create-menu {
  position: absolute;
  top: calc(100% + 12px);
  right: 0;
  min-width: 200px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  z-index: 2000;
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

/* 操作按钮 */
.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 12px;
  border: none;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  font-weight: 500;
}

.action-btn:hover {
  background: #667eea;
  color: white;
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
