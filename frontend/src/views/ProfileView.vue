<template>
  <div class="profile-view">
    <NavBar />

    <div class="profile-container">
      <!-- 左右布局 -->
      <div class="profile-layout">
        <!-- 左侧：自我介绍 -->
        <div class="left-section">
          <div class="self-intro">
            <div class="intro-header">
              <img src="@/assets/my_headportrait.jpg" alt="Avatar" class="intro-avatar" />
              <div class="intro-title">
                <h1>山角函兽</h1>
                <p>此情可待成追忆 --Keranthos的来源</p>
              </div>
            </div>
            <div class="intro-content">
              <div class="intro-text" v-html="formattedIntro"></div>
              <!-- 底部装饰图片 -->
              <div class="intro-bottom-image">
                <img src="/images/sunset-mountains.jpg" alt="Sunset Mountains" />
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：作品集 -->
        <div class="right-section">
          <div class="portfolio">
            <h2 class="portfolio-title">作品集</h2>
            <div class="portfolio-grid">
              <div
                v-for="project in projects"
                :key="project.id"
                class="project-card"
                @click="openProject(project)"
              >
                <div
                  class="project-content"
                  @mouseenter="showProjectImage(project)"
                  @mouseleave="hideProjectImage"
                >
                  <div class="project-header">
                    <div class="project-title">{{ project.title }}</div>
                    <div class="project-links">
                      <a
                        v-for="link in project.links"
                        :key="link.type"
                        :href="link.url"
                        :target="link.external ? '_blank' : '_self'"
                        class="project-link"
                        @click.stop
                        @mouseenter="hideProjectImage"
                        @mouseleave="showProjectImage(project)"
                      >
                        <FontAwesomeIcon :icon="getLinkIcon(link.type)" />
                      </a>
                    </div>
                  </div>
                  <div class="project-description">
                    {{ project.description }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 项目图片预览 -->
      <div
        v-if="hoveredProject && hoveredProject.image"
        class="project-image-preview"
        :style="imagePreviewStyle"
      >
        <img :src="hoveredProject.image" :alt="hoveredProject.title" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'
import NavBar from '@/components/NavBar.vue'
import { getArticlesNum, getArticlesList } from '@/api/Articles/browse'
import { getMomentsNum, getMomentsList } from '@/api/Moments/browse'

const stats = ref({
  totalArticles: 0,
  totalViews: 0,
  totalTags: 0,
  totalDays: 0
})

// 自我介绍内容
const introContent = ref('')
const hoveredProject = ref(null)
const mousePosition = ref({ x: 0, y: 0 })

// 项目数据
const projects = ref([
  {
    id: 1,
    title: 'EmulatorDAO',
    description: '基于BrokerChain构建的DAO组织，实现空投、社区贡献激励机制和NFT认证系统',
    image: require('@/assets/emulator-dao-flowchart.png'), // 流程图图片
    links: [
      {
        type: 'github',
        url: 'https://github.com/Keranthos/BrokerWallet_DAO_backend',
        label: 'GitHub 仓库',
        external: true
      },
      {
        type: 'article',
        url: 'https://mp.weixin.qq.com/s/SbZPfXkvlWZ3DyMCIChmSw',
        label: '技术文章',
        external: true
      }
    ]
  }
])

// 格式化自我介绍内容
const formattedIntro = computed(() => {
  return introContent.value
    .replace(/^# (.*$)/gim, '<h1>$1</h1>')
    .replace(/^## (.*$)/gim, '<h2>$1</h2>')
    .replace(/^### (.*$)/gim, '<h3>$1</h3>')
    .replace(/^\* (.*$)/gim, '<li>$1</li>')
    .replace(/^- (.*$)/gim, '<li>$1</li>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`(.*?)`/g, '<code>$1</code>')
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>')
    .replace(/\n\n/g, '</p><p>')
    .replace(/^(?!<[h1-3]|<li|<p)/gm, '<p>')
    .replace(/(<li>.*<\/li>)/g, '<ul>$1</ul>')
    .replace(/<\/ul>\s*<ul>/g, '')
})

// 图片预览样式
const imagePreviewStyle = computed(() => {
  return {
    left: `${mousePosition.value.x + 20}px`,
    top: `${mousePosition.value.y - 100}px`
  }
})

// 加载统计数据
const loadStats = async () => {
  try {
    const [blogCount, projectCount, researchCount, momentCount] = await Promise.all([
      getArticlesNum('blog'),
      getArticlesNum('project'),
      getArticlesNum('research'),
      getMomentsNum()
    ])

    stats.value.totalArticles = blogCount.num + projectCount.num + researchCount.num + momentCount.num

    // 获取真实的标签数量和写作天数
    await loadRealStats()
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载真实统计数据
const loadRealStats = async () => {
  try {
    const types = ['blog', 'project', 'research', 'moment']
    const allResults = []
    const tagSet = new Set()
    const dateSet = new Set()

    // 加载所有类型的文章
    for (const type of types) {
      if (type === 'moment') {
        const response = await getMomentsList(1, 1000)
        const articles = response.data.map(item => ({
          tags: [],
          createdAt: item.CreatedAt
        }))
        allResults.push(...articles)
      } else {
        const response = await getArticlesList(type, 1, 1000)
        const articles = response.data.map(item => ({
          tags: item.tags || [],
          createdAt: item.CreatedAt
        }))
        allResults.push(...articles)
      }
    }

    // 统计标签和写作日期
    allResults.forEach(article => {
      // 统计标签
      if (article.tags && article.tags.length > 0) {
        article.tags.forEach(tag => tagSet.add(tag))
      }

      // 统计写作日期（只统计日期，不统计时间）
      if (article.createdAt) {
        const date = new Date(article.createdAt)
        const dateStr = `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`
        dateSet.add(dateStr)
      }
    })

    // 计算总阅读量（模拟数据，实际项目中应该从后端获取）
    const totalViews = allResults.length * (Math.floor(Math.random() * 50) + 10)

    stats.value.totalViews = totalViews
    stats.value.totalTags = tagSet.size
    stats.value.totalDays = dateSet.size
  } catch (error) {
    console.error('加载真实统计数据失败:', error)
    // 如果失败，使用模拟数据
    stats.value.totalViews = Math.floor(Math.random() * 10000) + 1000
    stats.value.totalTags = Math.floor(Math.random() * 50) + 20
    stats.value.totalDays = Math.floor(Math.random() * 365) + 100
  }
}

// 显示项目图片
const showProjectImage = (project) => {
  if (project.image) {
    hoveredProject.value = project
  }
}

// 隐藏项目图片
const hideProjectImage = () => {
  // 添加淡出动画
  const preview = document.querySelector('.project-image-preview')
  if (preview) {
    preview.style.animation = 'imagePreviewSlideOut 0.25s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards'
    setTimeout(() => {
      hoveredProject.value = null
      if (preview) {
        preview.style.animation = ''
      }
    }, 250)
  } else {
    hoveredProject.value = null
  }
}

// 打开项目
const openProject = (project) => {
  // 这里可以添加项目详情页面的跳转逻辑
  console.log('打开项目:', project.title)
}

// 获取链接图标
const getLinkIcon = (type) => {
  const icons = {
    // 主要平台
    github: ['fab', 'github'],
    gitlab: ['fab', 'gitlab'],
    bitbucket: ['fab', 'bitbucket'],
    wechat: ['fab', 'weixin'],
    weibo: ['fab', 'weibo'],
    qq: ['fab', 'qq'],
    bilibili: ['fab', 'bilibili'],
    youtube: ['fab', 'youtube'],
    twitter: ['fab', 'twitter'],
    facebook: ['fab', 'facebook'],
    instagram: ['fab', 'instagram'],
    linkedin: ['fab', 'linkedin'],

    // 内容类型
    article: ['fas', 'newspaper'],
    blog: ['fas', 'blog'],
    post: ['fas', 'pen-fancy'],

    // 基础功能
    demo: ['fas', 'external-link-alt'],
    docs: ['fas', 'book'],
    website: ['fas', 'globe'],
    download: ['fas', 'download'],
    email: ['fas', 'envelope'],
    phone: ['fas', 'phone'],
    location: ['fas', 'map-marker-alt'],

    // 常用图标
    star: ['fas', 'star'],
    heart: ['fas', 'heart'],
    share: ['fas', 'share-alt'],
    edit: ['fas', 'edit'],
    delete: ['fas', 'trash'],
    search: ['fas', 'search'],
    refresh: ['fas', 'sync'],
    settings: ['fas', 'cog'],
    user: ['fas', 'user'],
    home: ['fas', 'home'],
    info: ['fas', 'info-circle'],
    warning: ['fas', 'exclamation-triangle'],
    error: ['fas', 'times-circle'],
    success: ['fas', 'check-circle'],
    question: ['fas', 'question-circle'],
    contact: ['fas', 'address-book'],

    // 媒体相关
    image: ['fas', 'image'],
    video: ['fas', 'video'],
    audio: ['fas', 'volume-up'],
    file: ['fas', 'file'],
    folder: ['fas', 'folder'],

    // 开发相关
    code: ['fas', 'code'],
    terminal: ['fas', 'terminal'],
    database: ['fas', 'database'],
    server: ['fas', 'server'],
    cloud: ['fas', 'cloud'],

    // 版本控制
    git: ['fab', 'git'],
    svn: ['fas', 'code'],
    mercurial: ['fas', 'code'],
    cvs: ['fas', 'code'],
    perforce: ['fas', 'code'],
    tfs: ['fas', 'code'],
    vss: ['fas', 'code'],
    clearcase: ['fas', 'code'],
    rcs: ['fas', 'code'],
    sccs: ['fas', 'code']
  }
  return icons[type] || ['fas', 'link']
}

// 加载自我介绍内容
const loadIntro = async () => {
  try {
    // 尝试从public目录加载
    const response = await fetch('/data/self-introduction.txt')
    if (response.ok) {
      introContent.value = await response.text()
    } else {
      throw new Error('文件未找到')
    }
  } catch (error) {
    console.error('加载自我介绍失败:', error)
    // 使用默认内容
    introContent.value = `我一直认为所有东西的价值都在于我已经或者即将赋予它的时间，山角函兽这个名字由来并没有引经据典、附庸风雅的故事：实际上仅仅是和三角函数谐音。
<br>
它可以令我时时想起最宝贵的一段时间，那时候拥有所有的少年意气，遇见了很多改变一生的人和事，乃至我一直希望它可以成为我人生的底色：永远充满希望、永远不会迷茫。同样它也象征着我另一项对人生的看法:衣不如新，人不如故。我所做的只是在某一天某一瞬间需要起一个花名的时候，从记忆的角落里把它翻了出来。
<br>
至于Keranthos就更加简单,单纯是我在需要一个英文花名时用AI把尝试翻译山角函兽后挑了一个最顺眼的名字:在某些时候我也是一名宿命论者,就是它了。`
  }
}

// 鼠标移动事件
const handleMouseMove = (event) => {
  mousePosition.value = {
    x: event.clientX,
    y: event.clientY
  }
}

onMounted(() => {
  loadStats()
  loadIntro()

  // 添加鼠标移动监听
  document.addEventListener('mousemove', handleMouseMove)
})

// 组件卸载时移除监听
onUnmounted(() => {
  document.removeEventListener('mousemove', handleMouseMove)
})
</script>

<style scoped>
.profile-view {
  min-height: 100vh;
  padding-top: 80px;
  position: relative;
}

.profile-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  position: relative;
  z-index: 1;
}

/* 左右布局 */
.profile-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  align-items: start;
}

/* 左侧：自我介绍 */
.left-section {
  position: relative;
}

.self-intro {
  background: transparent;
  padding: 30px;
}

.intro-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.intro-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 3px solid rgba(0, 0, 0, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.intro-title h1 {
  font-size: 2rem;
  color: #333;
  margin: 0 0 8px 0;
  font-weight: 600;
  text-align: left;
}

.intro-title p {
  font-size: 1rem;
  color: #666;
  margin: 0;
  text-align: left;
}

.intro-content {
  color: #333;
  line-height: 1.6;
  text-align: left;
  font-size: 1.1rem;
}

.intro-text h1 {
  font-size: 1.8rem;
  color: #333;
  margin: 20px 0 15px 0;
  font-weight: 600;
}

.intro-text h2 {
  font-size: 1.4rem;
  color: #333;
  margin: 18px 0 12px 0;
  font-weight: 500;
}

.intro-text h3 {
  font-size: 1.2rem;
  color: #333;
  margin: 15px 0 10px 0;
  font-weight: 500;
}

.intro-text p {
  margin: 12px 0;
  color: #333;
}

.intro-text ul {
  margin: 12px 0;
  padding-left: 20px;
}

.intro-text li {
  margin: 8px 0;
  color: #333;
}

.intro-text strong {
  color: #333;
  font-weight: 600;
}

.intro-text em {
  color: #666;
  font-style: italic;
}

.intro-text code {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  color: #333;
}

.intro-text a {
  color: #87ceeb;
  text-decoration: none;
  border-bottom: 1px solid rgba(135, 206, 235, 0.3);
  transition: all 0.3s ease;
}

.intro-text a:hover {
  color: #b0e0e6;
  border-bottom-color: rgba(176, 224, 230, 0.6);
}

/* 底部装饰图片样式 */
.intro-bottom-image {
  margin-top: 30px;
  text-align: center;
}

.intro-bottom-image img {
  max-width: 100%;
  height: auto;
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  opacity: 0.9;
  transition: all 0.3s ease;
}

.intro-bottom-image img:hover {
  opacity: 1;
  transform: translateY(-2px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
}

/* 右侧：作品集 */
.right-section {
  position: relative;
}

.portfolio {
  background: transparent;
  padding: 30px;
}

.portfolio-title {
  font-size: 1.8rem;
  color: #333;
  margin: 0 0 25px 0;
  font-weight: 600;
  text-align: left;
  position: relative;
}

.portfolio-title::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 2px;
}

.portfolio-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
}

.project-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(5px);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
  transition: all 0.3s ease;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.project-card:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.project-title {
  font-size: 1.2rem;
  color: #333;
  font-weight: 600;
  margin: 0;
  text-align: left;
}

.project-description {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.5;
  margin-bottom: 15px;
  text-align: left;
}

.project-links {
  display: flex;
  gap: 6px;
  align-items: center;
}

.project-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #f5f5f5;
  border-radius: 6px;
  color: #333;
  text-decoration: none;
  font-size: 1rem;
  transition: all 0.3s ease;
  border: 1px solid #ddd;
  margin-left: 8px;
}

.project-link:hover {
  background: #e0e0e0;
  color: #333;
  transform: translateY(-1px);
}

.project-link i {
  font-size: 0.9rem;
}

/* 链接按钮悬停效果 */
.project-link {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.project-link::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.project-link:hover {
  transform: translateY(-2px) scale(1.1);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.project-link:hover::before {
  left: 100%;
}

.project-link:active {
  transform: translateY(0) scale(1.05);
}

/* 项目图片预览 */
.project-image-preview {
  position: fixed;
  z-index: 1000;
  pointer-events: none;
  background: transparent;
  border-radius: 12px;
  padding: 0;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.2);
  max-width: 320px;
  max-height: 240px;
  overflow: hidden;
  transform: scale(0.8) translateY(10px);
  opacity: 0;
  animation: imagePreviewSlideIn 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
}

@keyframes imagePreviewSlideIn {
  0% {
    transform: scale(0.8) translateY(10px);
    opacity: 0;
  }
  100% {
    transform: scale(1) translateY(0);
    opacity: 1;
  }
}

@keyframes imagePreviewSlideOut {
  0% {
    transform: scale(1) translateY(0);
    opacity: 1;
  }
  100% {
    transform: scale(0.8) translateY(10px);
    opacity: 0;
  }
}

.project-image-preview img {
  width: 100%;
  height: auto;
  max-height: 240px;
  object-fit: contain;
  border-radius: 12px;
  background: transparent;
  display: block;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .profile-layout {
    grid-template-columns: 1fr;
    gap: 30px;
  }

  .portfolio-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 20px 15px;
  }

  .intro-header {
    flex-direction: column;
    text-align: left;
    gap: 15px;
    align-items: flex-start; /* 移动端纵向布局时，头像与标题块左对齐 */
  }

  .intro-title h1 {
    font-size: 1.6rem;
  }

  .portfolio-grid {
    grid-template-columns: 1fr;
  }

  .project-links {
    justify-content: flex-start;
  }
}

@media (max-width: 480px) {
  .self-intro,
  .portfolio {
    padding: 20px;
  }

  .intro-title h1 {
    font-size: 1.4rem;
  }

  .project-header {
    flex-direction: column;
    text-align: left;
    gap: 10px;
    align-items: flex-start;
  }

  .intro-bottom-image {
    margin-top: 20px;
  }

  .intro-bottom-image img {
    border-radius: 8px;
  }
}
</style>
