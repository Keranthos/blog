<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="detail-view">
    <!-- 文章标题区域 -->
    <div class="article-title-section">
      <div class="title-container">
        <h1 class="article-title">{{ title }}</h1>

        <!-- 标签和统计信息在同一行 -->
        <div class="tags-stats-row">
          <div class="article-tags">
            <span class="publish-date">{{ time }}</span>
            <span v-for="tag in tags" :key="tag" class="tag-item">{{ tag }}</span>
          </div>

          <div class="article-stats">
            <span class="word-count">本文字数 {{ content.length }} - 阅读时间约为 {{ readingTime ? readingTime.display : '未知' }}</span>
          </div>
        </div>

        <!-- 分割线 -->
        <div class="divider"></div>

        <!-- 四个互动按钮 -->
        <div class="engagement-buttons">
          <div class="left-buttons">
            <button class="like-btn">
              <font-awesome-icon icon="heart" />
              <span>0</span>
            </button>
            <button class="subscribe-btn">
              <font-awesome-icon icon="bookmark" />
              <span>订阅</span>
            </button>
          </div>

          <div class="right-buttons">
            <button class="comment-btn">
              <font-awesome-icon icon="comment" />
              <span>{{ comments.length }}</span>
            </button>
            <button class="share-btn" @click="toggleShareMenu">
              <font-awesome-icon icon="share" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 文章头图 -->
    <div class="article-image-section">
      <img :src="image" alt="Detail Image" class="article-image" loading="lazy" />
    </div>

    <!-- 文章内容卡片 -->
    <div ref="contentContainer" class="content-container">
      <!-- 标题位置占位容器 -->
      <div class="title-spacer"></div>

      <article class="article-content">
        <div class="markdown-body" v-html="renderedContent"></div>

        <!-- 悬浮文章目录 -->
        <div v-if="toc.length > 0 && tocVisible" class="floating-toc">
          <div class="toc-content">
            <div class="toc-header">
              <span class="toc-title">目录</span>
              <span class="toc-count">{{ toc.length }}</span>
            </div>
            <nav class="toc-nav">
              <div class="toc-tree">
                <div
                  v-for="(item, index) in toc"
                  :key="item.id"
                  :class="`toc-node toc-level-${item.relativeLevel}`"
                  :data-level="item.relativeLevel"
                  :data-index="index"
                >
                  <div class="toc-node-content">
                    <div class="toc-connectors">
                      <div
                        v-for="connector in item.connectors"
                        :key="`${connector.type}-${connector.level}`"
                        :class="`toc-connector toc-connector-${connector.type}`"
                        :style="connector.style"
                      ></div>
                    </div>
                    <a
                      :href="`#${item.id}`"
                      :title="item.text"
                      class="toc-link"
                      @click="scrollToHeading(item.id)"
                    >
                      <span class="toc-text">{{ item.text }}</span>
                    </a>
                  </div>
                </div>
              </div>
            </nav>
          </div>
        </div>
      </article>

      <!-- 分享菜单 -->
      <div v-if="showShareMenu" class="share-menu" @click.stop>
        <div class="share-platforms">
          <button class="share-platform-btn weibo" title="分享到微博" @click="shareToSocial('weibo')">
            <span class="platform-icon">📱</span>
            <span>微博</span>
          </button>
          <button class="share-platform-btn qq" title="分享到QQ" @click="shareToSocial('qq')">
            <span class="platform-icon">💬</span>
            <span>QQ</span>
          </button>
          <button class="share-platform-btn wechat" title="分享到微信" @click="shareToSocial('wechat')">
            <span class="platform-icon">💚</span>
            <span>微信</span>
          </button>
          <button class="share-platform-btn twitter" title="分享到Twitter" @click="shareToSocial('twitter')">
            <span class="platform-icon">🐦</span>
            <span>Twitter</span>
          </button>
          <button class="share-platform-btn facebook" title="分享到Facebook" @click="shareToSocial('facebook')">
            <span class="platform-icon">📘</span>
            <span>Facebook</span>
          </button>
          <button class="share-platform-btn copy" title="复制链接" @click="copyToClipboard()">
            <span class="platform-icon">🔗</span>
            <span>复制链接</span>
          </button>
        </div>
      </div>

      <!-- 相关文章推荐 -->
      <RelatedArticles
        v-if="articleData"
        :current-article="articleData"
        :max-count="3"
      />

      <!-- 评论区 -->
      <section class="comments-section">
        <div class="section-header">
          <h2>
            <font-awesome-icon icon="comment" class="section-icon" />
            评论区
            <span v-if="comments.length" class="comment-count">({{ comments.length }})</span>
          </h2>
        </div>

        <!-- 评论列表 -->
        <div v-if="Array.isArray(comments) && comments.length > 0" class="comments-list">
          <div v-for="comment in comments" :key="comment.ID" class="comment-item">
            <div class="comment-avatar">
              <div class="avatar-circle">{{ comment.username.charAt(0) }}</div>
            </div>
            <div class="comment-body">
              <div class="comment-header">
                <span class="comment-author">{{ comment.username }}</span>
                <span class="comment-time">{{ formatCommentTime(comment.CreatedAt) }}</span>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
            </div>
            <button v-if="user.level >= 3" class="delete-btn" @click="handleDeleteComment(comment.ID)">
              <font-awesome-icon icon="trash" />
            </button>
          </div>
        </div>

        <!-- 无评论提示 -->
        <div v-else class="empty-comments">
          <font-awesome-icon icon="comment-slash" class="empty-icon" />
          <p>还没有评论，来抢沙发吧！</p>
        </div>

        <!-- 发表评论 -->
        <div class="comment-editor">
          <h3>发表评论</h3>
          <div class="editor-wrapper">
            <textarea
              v-model="newComment"
              placeholder="说说你的想法..."
              class="comment-input"
              @keydown.ctrl.enter="submitComment"
            ></textarea>
            <div class="editor-actions">
              <span class="tip">Ctrl + Enter 快速发布</span>
              <button class="submit-btn" :disabled="!newComment.trim()" @click="submitComment">
                <font-awesome-icon icon="paper-plane" />
                发布评论
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { getArticleByID } from '@/api/Articles/browse'
import { getCommentsByID } from '@/api/Comments/browse'
import { createComment, deleteComment as deleteCommentAPI } from '@/api/Comments/edit'
import { showErrorMessage, showSuccessMessage, showWarningMessage } from '@/utils/waifuMessage'
import { marked } from 'marked'
import hljs from 'highlight.js'
import '@/assets/styles/github-highlight.css'
import { generateArticleSEO } from '@/utils/seo'
import RelatedArticles from '@/components/RelatedArticles.vue'
import { estimateReadingTime } from '@/utils/readingTime'

// 防抖函数
let scrollTimeout = null
const debounceScroll = (func, delay = 16) => {
  return (...args) => {
    clearTimeout(scrollTimeout)
    scrollTimeout = setTimeout(() => func.apply(this, args), delay)
  }
}

const props = defineProps({
  type: String, // 'blog', 'project', 'research', 'moment'
  articleId: String
})

const store = useStore()
const route = useRoute()
const user = store.state.user

// 添加复制按钮到代码块
const addCopyButton = (codeBlock) => {
  const pre = codeBlock.parentElement
  if (pre.querySelector('.copy-btn')) return // 避免重复添加

  // 等待pre元素完全渲染
  const checkAndAddButton = () => {
    if (!pre.offsetHeight || !pre.offsetWidth) {
      // 如果pre元素还没有尺寸，继续等待
      setTimeout(checkAndAddButton, 50)
      return
    }

    const copyBtn = document.createElement('button')
    copyBtn.className = 'copy-btn'
    copyBtn.innerHTML = '📋'
    copyBtn.title = '复制代码'

    // 确保复制按钮的定位样式
    copyBtn.style.position = 'absolute'
    copyBtn.style.top = '8px'
    copyBtn.style.right = '8px'
    copyBtn.style.left = 'auto'
    copyBtn.style.bottom = 'auto'
    copyBtn.style.zIndex = '999'

    copyBtn.addEventListener('click', async () => {
      try {
        await navigator.clipboard.writeText(codeBlock.textContent)
        showSuccessMessage('copy')
      } catch (err) {
        // 降级方案
        const textArea = document.createElement('textarea')
        textArea.value = codeBlock.textContent
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
        showSuccessMessage('copy')
      }
    })

    // 确保pre元素有正确的定位上下文
    pre.style.position = 'relative'
    pre.style.overflow = 'hidden' // 确保按钮不会超出边界
    pre.appendChild(copyBtn)

    // 调试信息已注释
    // console.log('复制按钮已添加到pre元素:', {
    //   preElement: pre,
    //   copyBtn,
    //   prePosition: pre.style.position,
    //   copyBtnPosition: copyBtn.style.position,
    //   preRect: pre.getBoundingClientRect()
    // })
  }

  // 开始检查并添加按钮
  checkAndAddButton()
}

// 精确调整目录位置 - 确保固定在屏幕右侧
const adjustTocPosition = () => {
  nextTick(() => {
    const tocElement = document.querySelector('.floating-toc')
    const contentContainer = document.querySelector('.content-container')

    if (!tocElement || !contentContainer) return

    // 清除任何可能冲突的left属性
    tocElement.style.left = 'auto'

    // 确保目录固定在屏幕右侧
    tocElement.style.right = '20px'
    tocElement.style.top = '80px'
    tocElement.style.position = 'fixed'

    // 位置调整完成后显示目录
    tocElement.style.opacity = '1'
    tocElement.style.pointerEvents = 'auto'
  })
}

// 生成文章目录 - 树形结构算法
const generateTOC = () => {
  const headings = document.querySelectorAll('.article-content h1, .article-content h2, .article-content h3, .article-content h4, .article-content h5, .article-content h6')
  const tocItems = []

  // 检测最高级标题
  let minLevel = 6
  headings.forEach((heading) => {
    const level = parseInt(heading.tagName.charAt(1))
    if (level < minLevel) {
      minLevel = level
    }
  })

  // 构建树形结构
  const treeNodes = []
  headings.forEach((heading, index) => {
    const level = parseInt(heading.tagName.charAt(1))
    const text = heading.textContent.trim()
    const id = `heading-${index}`

    // 为标题添加 ID
    heading.id = id

    // 计算相对层级
    const relativeLevel = level - minLevel + 1

    const node = {
      id,
      text,
      level,
      relativeLevel,
      element: heading,
      index,
      children: [],
      parent: null
    }

    treeNodes.push(node)
  })

  // 构建树形关系
  const rootNodes = []
  const stack = []

  treeNodes.forEach((node, index) => {
    // 弹出栈中比当前层级高的节点
    while (stack.length > 0 && stack[stack.length - 1].level >= node.level) {
      stack.pop()
    }

    // 设置父子关系
    if (stack.length > 0) {
      const parent = stack[stack.length - 1]
      node.parent = parent
      parent.children.push(node)
    } else {
      rootNodes.push(node)
    }

    stack.push(node)
  })

  // 为每个节点生成连接线
  treeNodes.forEach((node, index) => {
    const connectors = []

    // 为非顶级节点添加水平连接线
    if (node.relativeLevel > 1) {
      const lineLeft = 3 + (node.relativeLevel - 2) * 15
      const lineWidth = 8

      connectors.push({
        type: 'horizontal',
        level: node.relativeLevel,
        style: {
          left: `${lineLeft}px`,
          top: '50%',
          width: `${lineWidth}px`,
          height: '1px',
          position: 'absolute',
          transform: 'translateY(-50%)'
        }
      })
    }

    // 检查是否需要继续的垂直线（如果当前节点有兄弟节点）
    const hasSiblings = treeNodes.some((n, i) =>
      i > index &&
      n.level === node.level &&
      n.parent === node.parent
    )

    if (hasSiblings && node.relativeLevel > 1) {
      // 从当前节点继续向下的垂直线
      const verticalLeft = 3 + (node.relativeLevel - 2) * 15
      connectors.push({
        type: 'vertical-continue',
        level: node.relativeLevel,
        style: {
          left: `${verticalLeft}px`,
          top: '50%',
          width: '1px',
          height: '50%',
          position: 'absolute'
        }
      })
    }

    // 检查是否需要从父级继续的垂直线
    if (node.parent) {
      const parentLevel = node.parent.relativeLevel
      const currentLevel = node.relativeLevel

      // 从父级层级到当前层级的垂直线
      for (let level = parentLevel; level < currentLevel; level++) {
        const verticalLeft = 3 + (level - 1) * 15
        connectors.push({
          type: 'vertical',
          level,
          style: {
            left: `${verticalLeft}px`,
            top: '0px',
            width: '1px',
            height: '50%',
            position: 'absolute'
          }
        })
      }
    }

    tocItems.push({
      id: node.id,
      text: node.text,
      level: node.level,
      relativeLevel: node.relativeLevel,
      connectors,
      element: node.element,
      nodeIndex: node.index
    })
  })

  toc.value = tocItems

  // 生成目录后调整位置
  adjustTocPosition()

  // console.log('Generated TOC with tree structure:', {
  //   totalItems: toc.value.length,
  //   minLevel,
  //   treeStructure: toc.value.map(item => ({
  //     text: item.text.substring(0, 20) + '...',
  //     level: item.level,
  //     relativeLevel: item.relativeLevel,
  //     connectors: item.connectors.length,
  //     connectorTypes: item.connectors.map(c => c.type)
  //   }))
  // })

  // 延迟检查DOM中的连接线
  nextTick(() => {
    // const connectorElements = document.querySelectorAll('.toc-connector')
    // console.log('DOM中的连接线元素数量:', connectorElements.length)
    // connectorElements.forEach((el, index) => {
    //   const rect = el.getBoundingClientRect()
    //   console.log(`连接线 ${index}:`, {
    //     type: el.className,
    //     visible: rect.width > 0 && rect.height > 0,
    //     position: {
    //       left: el.style.left,
    //       top: el.style.top,
    //       width: el.style.width,
    //       height: el.style.height
    //     }
    //   })
    // })
  })
}

// 滚动到指定标题
const scrollToHeading = (id) => {
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

// 分享功能
const toggleShareMenu = () => {
  showShareMenu.value = !showShareMenu.value
}

// 切换目录显示功能已移除

// 智能目录跟随滚动逻辑
const handleScroll = () => {
  if (!contentContainer.value) return

  const scrollTop = window.scrollY
  const containerRect = contentContainer.value.getBoundingClientRect()
  const containerTop = containerRect.top + scrollTop

  // 确保目录在可见区域内才显示
  const containerBottom = containerRect.bottom + scrollTop
  tocVisible.value = scrollTop >= containerTop - 200 &&
                     scrollTop <= containerBottom - 200

  // 智能目录跟随滚动 - 使用requestAnimationFrame确保流畅
  requestAnimationFrame(() => {
    updateTocScrollPosition()
  })
}

// 更新目录滚动位置，确保当前标题在可视区域内
const updateTocScrollPosition = () => {
  if (!toc.value.length) return

  const tocNav = document.querySelector('.toc-nav')
  if (!tocNav) return

  // 获取当前视口中心位置
  const viewportCenter = window.innerHeight / 2
  const scrollTop = window.scrollY
  const centerY = scrollTop + viewportCenter

  // 找到最接近视口中心的标题
  let closestHeading = null
  let minDistance = Infinity

  toc.value.forEach((item, index) => {
    const element = document.getElementById(item.id)
    if (element) {
      const rect = element.getBoundingClientRect()
      const elementCenter = rect.top + scrollTop + rect.height / 2
      const distance = Math.abs(elementCenter - centerY)

      if (distance < minDistance) {
        minDistance = distance
        closestHeading = { item, index, element }
      }
    }
  })

  if (closestHeading) {
    // 找到对应的目录项元素
    const tocLinks = document.querySelectorAll('.toc-link')
    const targetTocLink = Array.from(tocLinks).find(link =>
      link.getAttribute('href') === `#${closestHeading.item.id}`
    )

    if (targetTocLink) {
      const tocVisibleHeight = tocNav.clientHeight
      const currentScrollTop = tocNav.scrollTop

      // 获取目录项的实际位置和尺寸
      const tocItemRect = targetTocLink.getBoundingClientRect()
      const tocNavRect = tocNav.getBoundingClientRect()

      // 计算目录项相对于目录容器的位置
      const itemTop = tocItemRect.top - tocNavRect.top + currentScrollTop
      const itemBottom = itemTop + tocItemRect.height

      // 检查当前标题是否在可视区域内
      const isItemVisible = itemTop >= currentScrollTop &&
                           itemBottom <= currentScrollTop + tocVisibleHeight

      // 如果不在可视区域内，则滚动到合适位置
      if (!isItemVisible) {
        let targetScrollTop

        // 如果标题在可视区域上方，滚动到标题位置
        if (itemTop < currentScrollTop) {
          targetScrollTop = itemTop
        } else if (itemBottom > currentScrollTop + tocVisibleHeight) {
          // 如果标题在可视区域下方，滚动到标题底部
          targetScrollTop = itemBottom - tocVisibleHeight
        }

        // 处理边界值
        const maxScrollTop = Math.max(0, tocNav.scrollHeight - tocVisibleHeight)
        const finalScrollTop = Math.max(0, Math.min(targetScrollTop, maxScrollTop))

        // 立即滚动到目标位置，确保目录跟随
        tocNav.scrollTo({
          top: finalScrollTop,
          behavior: 'instant'
        })
      }
    }

    // 更新当前活跃的目录项
    updateActiveTocItem(closestHeading.item.id)
  }
}

// 更新当前活跃的目录项
const updateActiveTocItem = (activeId) => {
  const tocLinks = document.querySelectorAll('.toc-link')
  tocLinks.forEach(link => {
    link.classList.remove('active')
    if (link.getAttribute('href') === `#${activeId}`) {
      link.classList.add('active')
    }
  })
}

// 初始化目录滚动位置，确保第一个标题可见
const initializeTocScrollPosition = () => {
  if (!toc.value.length) return

  const tocNav = document.querySelector('.toc-nav')
  if (!tocNav) return

  // 将目录滚动到顶部，确保第一个标题可见
  tocNav.scrollTo({
    top: 0,
    behavior: 'instant'
  })

  // 激活第一个标题
  if (toc.value.length > 0) {
    updateActiveTocItem(toc.value[0].id)
  }
}

const shareToSocial = (platform) => {
  const url = window.location.href
  const shareText = `${title.value} - ${window.location.origin}`

  let shareUrl = ''

  switch (platform) {
    case 'weibo':
      shareUrl = `https://service.weibo.com/share/share.php?url=${encodeURIComponent(url)}&title=${encodeURIComponent(shareText)}`
      break
    case 'qq':
      shareUrl = `https://connect.qq.com/widget/shareqq/index.html?url=${encodeURIComponent(url)}&title=${encodeURIComponent(shareText)}`
      break
    case 'wechat':
      // 微信分享需要特殊处理，这里显示提示
      showWarningMessage('微信分享需要手动复制链接哦～')
      copyToClipboard()
      return
    case 'twitter':
      shareUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(shareText)}&url=${encodeURIComponent(url)}`
      break
    case 'facebook':
      shareUrl = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`
      break
  }

  if (shareUrl) {
    window.open(shareUrl, '_blank', 'width=600,height=400')
    showShareMenu.value = false
  }
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(window.location.href)
    showSuccessMessage('copy')
    showShareMenu.value = false
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = window.location.href
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showSuccessMessage('copy')
    showShareMenu.value = false
  }
}

// 编辑文章功能已移除，因为新的布局中没有编辑按钮
const image = ref('')
const title = ref('')
const tags = ref([])
const time = ref('')
const content = ref('')
const renderedContent = ref('')
const comments = ref([])
const newComment = ref('')
const viewCount = ref(0)
const toc = ref([]) // 文章目录
const showShareMenu = ref(false) // 显示分享菜单
const readingTime = ref(null) // 阅读时间估算
const tocVisible = ref(true) // 目录是否可见
const contentContainer = ref(null) // 内容容器引用

// 文章数据对象（用于相关文章推荐）
const articleData = computed(() => {
  if (!title.value) return null

  return {
    id: props.articleId || route.params.id,
    type: props.type || route.params.type,
    title: title.value,
    content: content.value,
    image: image.value,
    tags: tags.value,
    time: time.value,
    viewCount: viewCount.value
  }
})

// 加载文章详情
const loadDetail = async () => {
  const id = props.articleId || route.params.id
  // console.log('ArticleDetail - 获取ID:', {
  //   propsArticleId: props.articleId,
  //   routeParamsId: route.params.id,
  //   finalId: id,
  //   route: route.path,
  //   routeParams: route.params
  // })

  // 如果ID为undefined，直接返回错误
  if (!id || id === 'undefined') {
    console.error('ArticleDetail - ID无效:', id)
    showErrorMessage('文章ID无效')
    return
  }

  let response

  try {
    if (props.type === 'moment') {
    // 碎碎念使用不同的 API
      const { getMoment } = await import('@/api/Moments/browse')
      const res = await getMoment(id)
      response = res.data
      // 碎碎念的字段映射
      image.value = response.Image || 'https://picsum.photos/id/180/1920/1080'
      title.value = response.Title
      tags.value = []
      time.value = response.CreatedAt ? new Date(response.CreatedAt).toLocaleDateString('zh-CN') : '未知时间'
      content.value = response.Content
      viewCount.value = response.viewCount || 0

      // 计算阅读时间
      readingTime.value = estimateReadingTime(response.Content, response.Title)
    } else {
      const res = await getArticleByID(props.type, id)
      response = res.data
      image.value = response.image
      title.value = response.title
      tags.value = response.tags
      time.value = response.time || response.CreatedAt ? new Date(response.time || response.CreatedAt).toLocaleDateString('zh-CN') : '未知时间'
      content.value = response.content
      viewCount.value = response.viewCount || 0

      // 计算阅读时间
      readingTime.value = estimateReadingTime(response.content, response.title)
    }
  } catch (error) {
    console.error('加载文章详情失败:', error)

    // 处理404错误
    if (error.response && error.response.status === 404) {
      // 设置404页面内容
      image.value = 'https://picsum.photos/id/180/1920/1080'
      title.value = '文章未找到'
      tags.value = []
      time.value = '未知时间'
      content.value = '抱歉，您访问的文章不存在或已被删除。'
      viewCount.value = 0
      readingTime.value = null

      // 显示错误消息
      showErrorMessage('文章不存在')
      return
    }

    // 处理其他错误
    showErrorMessage('加载文章失败，请稍后重试')
    return
  }

  // 渲染 Markdown 内容
  renderedContent.value = marked(content.value, {
    breaks: false, // 不将换行符转换为 <br>
    gfm: true, // 启用 GitHub Flavored Markdown
    headerIds: false,
    mangle: false,
    pedantic: false, // 禁用严格模式
    sanitize: false, // 不禁用HTML标签
    smartLists: true,
    smartypants: false
  })

  // 后处理：修复没有被正确渲染的粗体语法
  // 将 **text:** 这样的模式手动转换为 <strong>text:</strong>
  renderedContent.value = renderedContent.value.replace(/\*\*([^*:]+:\**)\*\*/g, '<strong>$1</strong>')

  // 更新SEO信息
  const articleData = {
    ID: id,
    title: title.value,
    content: content.value,
    abstract: response.abstract || '',
    tags: tags.value,
    image: image.value,
    CreatedAt: time.value,
    UpdatedAt: time.value
  }
  const seoData = generateArticleSEO(articleData, props.type || route.params.type)
  import('@/utils/seo').then(({ updateSEO }) => {
    updateSEO(seoData)
  })

  // 等待 DOM 更新后手动触发代码高亮和添加复制按钮
  nextTick(() => {
    // 先修复可能残留的 **粗体** 文本节点（不影响已正确渲染的 strong 标签）
    fixResidualBoldInDOM()

    // 延迟执行，确保代码块完全渲染
    setTimeout(() => {
      document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightElement(block)
        addCopyButton(block)
      })

      // 生成目录
      generateTOC()
    }, 100) // 延迟100ms确保渲染完成
  })
}

// 加载评论
const loadComments = async () => {
  try {
    const id = props.articleId || route.params.id

    // 如果ID无效，跳过评论加载
    if (!id || id === 'undefined') {
      // console.log('ArticleDetail - 跳过评论加载，ID无效:', id)
      comments.value = []
      return
    }

    const res = await getCommentsByID(props.type, id)
    comments.value = res.data
  } catch (error) {
    console.error('加载评论失败:', error)
    comments.value = []
    // 评论加载失败不影响页面显示，静默处理
  }
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) return

  const id = props.articleId || route.params.id
  const user = store.state.user

  if (!id || id === 'undefined') {
    showErrorMessage('文章ID无效，无法提交评论')
    return
  }

  if (!user || !user.token) {
    showErrorMessage('401')
    return
  }

  try {
    await createComment(user, id, props.type, newComment.value)
    newComment.value = ''
    await loadComments()
    showSuccessMessage('comment')
  } catch (error) {
    showErrorMessage(error)
  }
}

// 删除评论
const handleDeleteComment = async (commentId) => {
  // 第一次确认
  const firstConfirm = confirm('确定要删除这条评论吗？\n\n⚠️ 此操作不可撤销！')
  if (!firstConfirm) return

  // 第二次确认（防误触）
  const secondConfirm = confirm('再次确认：真的要删除这条评论吗？\n\n删除后将无法恢复！')
  if (!secondConfirm) return

  try {
    const user = store.state.user
    await deleteCommentAPI(user, commentId)
    comments.value = comments.value.filter((comment) => comment.ID !== commentId)

    // 显示看板娘消息（如果可用）
    if (window.showMessage) {
      window.showMessage('评论删除成功～', 3000, 9)
    }
  } catch (error) {
    console.error('删除评论失败:', error)

    // 显示看板娘错误消息（如果可用）
    if (window.showMessage) {
      window.showMessage('(｡•́︿•̀｡)<br>删除失败了…请重试吧～', 5000, 10)
    }
  }
}

// 在组件挂载时加载文章和评论
onMounted(async () => {
  // 重置页面滚动位置到顶部（使用smooth避免突然跳转触发hover）
  window.scrollTo({ top: 0, behavior: 'instant' })

  // 延迟执行，确保路由参数完全加载
  await nextTick()

  const id = props.articleId || route.params.id
  // console.log('ArticleDetail - 组件初始化:', {
  //   propsArticleId: props.articleId,
  //   routeParamsId: route.params.id,
  //   finalId: id
  // })

  // 只有当ID有效时才加载数据
  if (id && id !== 'undefined') {
    await loadDetail()
    await loadComments()
  } else {
    console.warn('ArticleDetail - 组件初始化时ID无效，跳过数据加载')
  }

  // 点击外部关闭分享菜单
  document.addEventListener('click', (e) => {
    if (!e.target.closest('.share-item')) {
      showShareMenu.value = false
    }
  })

  // 添加滚动监听器（使用防抖优化性能）
  const debouncedHandleScroll = debounceScroll(handleScroll, 16)
  window.addEventListener('scroll', debouncedHandleScroll, { passive: true })

  // 添加窗口大小变化监听器 - 延迟执行确保布局稳定
  window.addEventListener('resize', () => {
    setTimeout(adjustTocPosition, 100)
  }, { passive: true })

  // 初始化目录位置 - 延迟确保DOM完全渲染
  setTimeout(() => {
    adjustTocPosition()
    // 初始化目录滚动位置，确保第一个标题可见
    initializeTocScrollPosition()
  }, 300)

  // 初始化滚动状态
  nextTick(() => {
    handleScroll()
  })
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  // 清理滚动监听器
  if (scrollTimeout) {
    clearTimeout(scrollTimeout)
  }
  // 注意：由于使用了防抖函数，这里需要保存函数引用才能正确移除
  // 在实际应用中，如果需要严格清理，应该保存函数引用
})

// 格式化评论时间
const formatCommentTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`

  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 监听路由变化，重新加载文章和评论
watch(
  () => route.params.id,
  async (newId) => {
    // console.log('ArticleDetail - 路由变化:', { newId, routeParams: route.params })

    // 重置页面滚动位置到顶部（使用instant避免触发hover）
    window.scrollTo({ top: 0, behavior: 'instant' })

    // 只有当新ID存在且有效时才重新加载
    if (newId && newId !== 'undefined') {
      await loadDetail()
      await loadComments()
    }
  }
)

const fixResidualBoldInDOM = () => {
  const container = document.querySelector('.markdown-body')
  if (!container) return

  const isInCodeLike = (node) => {
    let current = node.parentNode
    while (current) {
      const tag = current.nodeName
      if (tag === 'CODE' || tag === 'PRE' || tag === 'KBD' || tag === 'SAMP') return true
      current = current.parentNode
    }
    return false
  }

  const walker = document.createTreeWalker(
    container,
    NodeFilter.SHOW_TEXT,
    {
      acceptNode: (node) => {
        if (!node.nodeValue || node.nodeValue.indexOf('**') === -1) {
          return NodeFilter.FILTER_REJECT
        }
        if (isInCodeLike(node)) return NodeFilter.FILTER_REJECT
        return NodeFilter.FILTER_ACCEPT
      }
    },
    false
  )

  const candidates = []
  let textNode = walker.nextNode()
  while (textNode) {
    candidates.push(textNode)
    textNode = walker.nextNode()
  }

  const boldPattern = /\*\*([^*]+?)\*\*/g

  candidates.forEach((node) => {
    const text = node.nodeValue
    if (!boldPattern.test(text)) return
    boldPattern.lastIndex = 0

    const fragment = document.createDocumentFragment()
    let lastIndex = 0
    let match

    while ((match = boldPattern.exec(text)) !== null) {
      const [full, inner] = match
      const start = match.index
      const end = start + full.length

      if (start > lastIndex) {
        fragment.appendChild(document.createTextNode(text.slice(lastIndex, start)))
      }

      const strong = document.createElement('strong')
      strong.textContent = inner
      fragment.appendChild(strong)

      lastIndex = end
    }

    if (lastIndex < text.length) {
      fragment.appendChild(document.createTextNode(text.slice(lastIndex)))
    }

    if (node.parentNode) {
      node.parentNode.replaceChild(fragment, node)
    }
  })
}
</script>

<style scoped>
.detail-view {
  min-height: 100vh;
  padding-top: 80px;
  background: transparent; /* 移除默认背景，使用全局背景 */
}

/* 文章标题区域 */
.article-title-section {
  padding: 40px 0 20px;
}

.title-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 40px;
}

.article-title {
  color: #333;
  font-size: 2.2rem;
  font-weight: 700;
  margin-bottom: 20px;
  line-height: 1.3;
  text-align: left;
}

/* 标签和统计信息在同一行 */
.tags-stats-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  flex-wrap: wrap;
  gap: 15px;
}

.article-tags {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.tag-item {
  display: inline-block;
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 0.8rem;
  border-radius: 12px;
  font-weight: 500;
}

.publish-date {
  font-size: 0.9rem;
  color: #666;
}

/* 分割线 */
.divider {
  height: 1px;
  background: #e0e0e0;
  margin-bottom: 15px;
}

/* 四个互动按钮 */
.engagement-buttons {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.left-buttons, .right-buttons {
  display: flex;
  gap: 15px;
}

.like-btn, .subscribe-btn, .comment-btn, .share-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  color: #666;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.like-btn:hover, .subscribe-btn:hover, .comment-btn:hover, .share-btn:hover {
  background: #f8f9fa;
  border-color: #d0d0d0;
  color: #333;
}

.like-btn svg {
  color: #ff6b6b;
}

.subscribe-btn svg {
  color: #ffa500;
}

.comment-btn svg {
  color: #4a90e2;
}

.share-btn svg {
  color: #50c878;
}

/* 文章统计信息 */
.article-stats {
  color: #999;
  font-size: 0.85rem;
}

/* 文章头图区域 */
.article-image-section {
  width: 100%;
  margin: 20px 0;
  text-align: center;
}

.article-image {
  max-width: 100%;
  max-height: 450px;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

/* 内容容器 */
.content-container {
  max-width: 1000px;
  margin: 0 auto 60px;
  position: relative;
  z-index: 2;
  padding: 0 40px;
}

/* 标题位置占位容器 - 已不需要，因为标题现在在顶部 */
.title-spacer {
  display: none;
}

/* 文章内容卡片 */
.article-content {
  position: relative; /* 确保目录以此为定位基准 */
  background: transparent; /* 移除白色背景，与整体背景融合 */
  border-radius: 0; /* 移除圆角 */
  padding: 50px 0; /* 只保留上下内边距 */
  box-shadow: none; /* 移除阴影 */
  margin-bottom: 40px;
}

.markdown-body {
  line-height: 1.8;
  color: #333;
  font-size: 1.05rem;
  text-align: left !important;
  background: transparent !important; /* 覆盖github-markdown.css中的白色背景 */
}

/* 文章内容中的图片样式 */
.markdown-body img {
  max-width: 80% !important;
  max-height: 400px !important;
  width: auto !important;
  height: auto !important;
  display: block !important;
  margin: 20px auto !important;
  border-radius: 8px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
  object-fit: contain !important;
}

.markdown-body * {
  text-align: left !important;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  text-align: left !important;
}

.markdown-body p,
.markdown-body div,
.markdown-body span {
  text-align: left !important;
}

/* 确保粗体样式正确显示 */
.markdown-body strong,
.markdown-body b {
  font-weight: 600 !important;
  color: inherit !important;
}

/* 强制确保所有strong标签都显示为粗体 */
.markdown-body strong {
  font-weight: bold !important;
  font-weight: 700 !important;
  font-family: 'Inter', 'Noto Sans SC', sans-serif !important;
}

/* 调试样式已移除 */

/* 评论区 */
.comments-section {
  background: rgba(255, 255, 255, 0.5);
  border-radius: 0;
  padding: 30px;
  box-shadow: none;
  border: none;
  margin-bottom: 30px;
  backdrop-filter: blur(25px);
}

.section-header {
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.section-header h2 {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.5rem;
  color: #333;
  margin: 0;
  font-weight: 600;
}

.section-icon {
  color: #667eea;
  font-size: 1.5rem;
}

.comment-count {
  font-size: 1rem;
  color: #999;
  font-weight: 400;
}

/* 评论列表 */
.comments-list {
  margin-bottom: 40px;
}

.comment-item {
  display: flex;
  gap: 15px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 0;
  margin-bottom: 12px;
  transition: all 0.3s ease;
  border: none;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  backdrop-filter: blur(15px);
}

.comment-item:hover {
  background: rgba(255, 255, 255, 0.18);
  border-bottom-color: rgba(102, 126, 234, 0.2);
  box-shadow: none;
  backdrop-filter: blur(20px);
}

.comment-avatar {
  flex-shrink: 0;
}

.avatar-circle {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: 600;
  text-transform: uppercase;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 600;
  color: #333;
  font-size: 1rem;
}

.comment-time {
  color: #999;
  font-size: 0.85rem;
}

.comment-content {
  color: #555;
  line-height: 1.6;
  word-wrap: break-word;
  font-size: 0.95rem;
}

.delete-btn {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.delete-btn:hover {
  background: #ff4d4f;
  color: white;
  transform: scale(1.1);
}

/* 无评论提示 */
.empty-comments {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-icon {
  font-size: 4rem;
  color: #ddd;
  margin-bottom: 20px;
}

.empty-comments p {
  font-size: 1.1rem;
  margin: 0;
}

/* 评论编辑器 */
.comment-editor {
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
  background: rgba(255, 255, 255, 0.15);
  border-radius: 0;
  padding: 30px;
  backdrop-filter: blur(25px);
}

.comment-editor h3 {
  font-size: 1.3rem;
  color: #333;
  margin-bottom: 20px;
}

.editor-wrapper {
  background: rgba(255, 255, 255, 0.12);
  border-radius: 0;
  padding: 20px;
  border: none;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  transition: all 0.3s ease;
  backdrop-filter: blur(15px);
}

.editor-wrapper:focus-within {
  border-bottom-color: rgba(102, 126, 234, 0.3);
  background: rgba(255, 255, 255, 0.18);
  box-shadow: none;
  backdrop-filter: blur(20px);
}

.comment-input {
  width: 100%;
  min-height: 120px;
  padding: 15px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  line-height: 1.6;
  resize: vertical;
  background: transparent;
  color: #333;
  font-family: inherit;
}

.comment-input:focus {
  outline: none;
}

.comment-input::placeholder {
  color: #aaa;
}

.editor-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
}

.tip {
  color: #999;
  font-size: 0.85rem;
}

/* 代码复制按钮样式 - 简洁版 */
.copy-btn {
  position: absolute !important;
  top: 8px !important;
  right: 8px !important;
  left: auto !important;
  bottom: auto !important;
  width: 24px !important;
  height: 24px !important;
  padding: 0 !important;
  margin: 0 !important;
  background: rgba(0, 0, 0, 0.6) !important;
  color: #fff !important;
  border: none !important;
  border-radius: 4px !important;
  font-size: 12px !important;
  cursor: pointer !important;
  transition: all 0.2s ease !important;
  backdrop-filter: blur(10px) !important;
  z-index: 999 !important;
  opacity: 0 !important;
  transform: scale(0.8) !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  line-height: 1 !important;
  text-align: center !important;
}

.copy-btn:hover {
  background: rgba(0, 0, 0, 0.8) !important;
  transform: scale(1) !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2) !important;
}

/* 代码块悬停时显示复制按钮 */
pre:hover .copy-btn {
  opacity: 1 !important;
  transform: scale(1) !important;
}

/* 确保pre元素有正确的定位上下文 */
.markdown-body pre {
  position: relative !important;
  overflow: hidden !important;
}

/* 文章目录样式 - 初始隐藏，等待JS定位 */
.floating-toc {
  position: fixed;
  top: 80px;
  right: 20px; /* 向屏幕右边靠近 */
  width: 200px;
  max-height: 70vh;
  z-index: 10;
  transition: all 0.3s ease;
  opacity: 1; /* 始终可见 */
  pointer-events: auto; /* 可交互 */
}

/* 确保在不同屏幕尺寸下都紧贴正文 - 完全由JS控制 */

/* 目录切换按钮已移除 */

/* 目录内容 - 淡色边框和毛玻璃效果 */
.toc-content {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  box-shadow: none;
  overflow: visible;
  backdrop-filter: blur(20px);
  padding: 16px;
}

/* 隐藏目录头部 */
.toc-header {
  display: none !important;
}

.toc-title {
  display: none !important;
}

.toc-count {
  display: none !important;
}

/* 目录导航 */
.toc-nav {
  max-height: 60vh; /* 限制最大高度为视口高度的60% */
  overflow-y: auto; /* 允许垂直滚动 */
  overflow-x: hidden; /* 隐藏水平滚动 */
}

/* 隐藏滚动条 */
.toc-nav::-webkit-scrollbar {
  display: none;
}

.toc-nav {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}

/* 滚动条样式已移除，使用隐藏滚动条 */

.toc-tree {
  position: relative;
  padding: 0;
  margin: 0;
}

.toc-node {
  position: relative;
  margin: 0;
}

.toc-node-content {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 24px;
}

.toc-connectors {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 3;
}

.toc-connector {
  position: absolute;
  background: #d1d5da;
  z-index: 3;
  pointer-events: none;
}

.toc-connector-vertical {
  width: 1px;
  background: #d1d5da;
  opacity: 1;
}

.toc-connector-horizontal {
  height: 1px;
  background: #d1d5da;
  opacity: 1;
}

.toc-connector-vertical-continue {
  width: 1px;
  background: #d1d5da;
  opacity: 1;
}

.toc-link {
  display: block;
  padding: 8px 12px;
  color: #586069;
  text-decoration: none;
  transition: all 0.3s ease;
  font-size: 13px;
  line-height: 1.4;
  border-left: 3px solid transparent;
  position: relative;
  border-radius: 0;
  margin: 2px 0;
  z-index: 2;
  background: transparent !important;
}

.toc-link:hover {
  background: transparent !important;
  color: #0366d6;
  border-left-color: #0366d6;
  transform: translateX(3px);
  box-shadow: none;
}

.toc-link.active {
  background: transparent !important;
  color: #495057;
  border-left-color: #0366d6;
  box-shadow: none;
  transform: translateX(2px);
  font-weight: 600;
}

.toc-text {
  display: block;
  overflow: visible;
  text-overflow: unset;
  white-space: normal;
  word-wrap: break-word;
}

/* 层级样式 - 透明背景 */
.toc-level-1 .toc-link {
  font-weight: 700;
  font-size: 14px;
  color: #333;
  border-left: 3px solid #0366d6;
  background: transparent !important;
  margin: 4px 0;
  padding: 10px 12px;
}

.toc-level-2 .toc-link {
  font-weight: 600;
  font-size: 13px;
  color: #555;
  background: transparent !important;
  padding-left: 24px;
  padding-right: 12px;
}

.toc-level-3 .toc-link {
  font-weight: 500;
  font-size: 12px;
  color: #666;
  background: transparent !important;
  padding-left: 36px;
  padding-right: 12px;
}

.toc-level-4 .toc-link {
  font-weight: 400;
  font-size: 12px;
  color: #777;
  background: transparent !important;
  padding-left: 48px;
  padding-right: 12px;
}

.toc-level-5 .toc-link {
  font-weight: 400;
  font-size: 11px;
  color: #888;
  background: transparent !important;
  padding-left: 60px;
  padding-right: 12px;
}

.toc-level-6 .toc-link {
  font-weight: 400;
  font-size: 11px;
  color: #888;
  background: transparent !important;
  padding-left: 72px;
  padding-right: 12px;
}

/* 目录层级缩进 - 动态生成，静态样式已移除 */

/* 五级和六级标题样式由动态CSS生成 */

/* 层级指示器已移除，由动态样式控制 */

/* 五级和六级标题指示器已移除，由动态样式控制 */

/* 树形结构连接线样式由动态CSS生成 */

/* 响应式设计 */
@media (max-width: 768px) {
  /* 移动端隐藏悬浮目录 */
  .floating-toc {
    display: none;
  }
}

/* 平板端优化 */
@media (min-width: 769px) and (max-width: 1024px) {
  .floating-toc {
    width: 140px;
    /* 定位完全由JS控制 */
  }
}

/* 大屏幕优化 - 完全由JS控制 */

/* 超大屏幕优化 */
@media (min-width: 1600px) {
  .floating-toc {
    width: 160px;
    /* 定位完全由JS控制 */
  }
}

/* 分享功能样式 */
.share-item {
  position: relative;
}

.share-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 20px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.share-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.4);
  transform: translateY(-1px);
}

.share-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  z-index: 1000;
  min-width: 200px;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.share-platforms {
  padding: 12px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.share-platform-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 8px;
  background: transparent;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.8rem;
  color: #666;
}

.share-platform-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: rgba(102, 126, 234, 0.3);
  transform: translateY(-2px);
}

.platform-icon {
  font-size: 1.2rem;
}

/* 不同平台的特色颜色 */
.share-platform-btn.weibo:hover { background: rgba(230, 22, 45, 0.1); border-color: rgba(230, 22, 45, 0.3); }
.share-platform-btn.qq:hover { background: rgba(18, 183, 245, 0.1); border-color: rgba(18, 183, 245, 0.3); }
.share-platform-btn.wechat:hover { background: rgba(9, 187, 7, 0.1); border-color: rgba(9, 187, 7, 0.3); }
.share-platform-btn.twitter:hover { background: rgba(29, 161, 242, 0.1); border-color: rgba(29, 161, 242, 0.3); }
.share-platform-btn.facebook:hover { background: rgba(24, 119, 242, 0.1); border-color: rgba(24, 119, 242, 0.3); }
.share-platform-btn.copy:hover { background: rgba(156, 39, 176, 0.1); border-color: rgba(156, 39, 176, 0.3); }

/* 响应式设计 */
@media (max-width: 768px) {
  .share-menu {
    right: -50px;
    min-width: 180px;
  }

  .share-platforms {
    grid-template-columns: 1fr 1fr 1fr;
    padding: 8px;
  }

  .share-platform-btn {
    padding: 8px 4px;
    font-size: 0.7rem;
  }

  .platform-icon {
    font-size: 1rem;
  }
}

.submit-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  box-shadow: none;
}

/* 响应式 */
@media (max-width: 768px) {
  .article-title-section {
    padding: 20px 0 15px;
  }

  .title-container {
    padding: 0 20px;
  }

  .article-title {
    font-size: 1.8rem;
    margin-bottom: 15px;
  }

  .tags-stats-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 15px;
  }

  .article-tags {
    gap: 10px;
  }

  .tag-item {
    font-size: 0.75rem;
    padding: 3px 10px;
  }

  .engagement-buttons {
    flex-direction: column;
    gap: 15px;
    margin-bottom: 15px;
  }

  .left-buttons, .right-buttons {
    gap: 10px;
  }

  .like-btn, .subscribe-btn, .comment-btn, .share-btn {
    padding: 6px 12px;
    font-size: 0.8rem;
  }

  .article-stats {
    font-size: 0.8rem;
  }

  .article-image {
    max-height: 300px;
  }

  .content-container {
    padding: 0 20px;
  }

  .article-content {
    padding: 30px 20px;
  }

  .comments-section {
    padding: 25px 20px;
  }

  .comment-item {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
