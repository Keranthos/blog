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

        <!-- 四个互动按钮和编辑按钮 -->
        <div class="engagement-buttons">
          <div class="left-buttons">
            <button class="like-btn" :class="{ liked: isLiked, loading: likeLoading }" @click="handleLike" :disabled="likeLoading">
              <font-awesome-icon :icon="isLiked ? 'heart' : ['far', 'heart']" />
              <span>{{ likeCount }}</span>
            </button>
            <button class="subscribe-btn" @click="handleSubscribe">
              <font-awesome-icon icon="bookmark" />
              <span>订阅</span>
            </button>
          </div>

          <div class="right-buttons">
            <button v-if="user.level >= 3" class="edit-btn" @click="goToEdit">
              <font-awesome-icon icon="pen-to-square" />
              <span>编辑</span>
            </button>
            <button class="comment-btn" @click="scrollToComments">
              <font-awesome-icon icon="comment" />
              <span>{{ comments.length }}</span>
            </button>
            <button class="share-btn" @click="handleShare">
              <font-awesome-icon icon="share" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 文章头图 -->
    <div class="article-image-section">
      <img :src="image" alt="Detail Image" class="article-image" loading="lazy" decoding="async" @error="onImgError($event)" />
    </div>

    <!-- 文章内容卡片 -->
    <div ref="contentContainer" class="content-container">
      <!-- 标题位置占位容器 -->
      <div class="title-spacer"></div>

      <article ref="articleContentRef" class="article-content">
        <div class="markdown-body" v-html="renderedContent"></div>

        <!-- 文本选择菜单 -->
        <TextSelectionMenu
          :visible="textSelectionMenuVisible"
          :position="textSelectionPosition"
          :selected-text="selectedText"
          :is-highlighted="isTextHighlighted"
          @copy="handleTextCopy"
          @highlight="handleTextHighlight"
          @share="handleTextShare"
          @comment="handleTextComment"
        />

        <!-- 分享卡片 -->
        <ShareCard
          :visible="shareCardVisible"
          :selected-text="shareSelectedText"
          :article-title="title"
          :article-subtitle="getArticleSubtitle()"
          :article-url="articleUrl"
          @close="shareCardVisible = false"
        />

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
            <svg class="comment-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
            </svg>
            评论
            <span v-if="comments.length" class="comment-count">({{ comments.length }})</span>
          </h2>
        </div>

        <!-- 发表评论 -->
        <div class="comment-editor">
          <div class="editor-wrapper">
            <div class="comment-input-wrapper">
              <div class="comment-input-box" :class="{ 'has-quoted-text': quotedText }">
                <!-- 引用原文显示（在输入框内部） -->
                <div v-if="quotedText" class="quoted-text-container">
                  <div class="quoted-text-label">引用原文：</div>
                  <div
                    ref="inputQuotedTextRef"
                    class="quoted-text-content markdown-body"
                    :class="{ 'is-truncated': isInputQuotedTextTruncated }"
                    v-html="renderedQuotedText"
                  ></div>
                  <button class="cancel-quote-btn" @click="clearQuotedText" title="取消引用">
                    <font-awesome-icon icon="times" />
                  </button>
                </div>
                <textarea
                  ref="commentTextarea"
                  v-model="newComment"
                  placeholder="支持 Markdown 语法"
                  class="comment-input"
                  @keydown.ctrl.enter="submitComment"
                  @focus="handleInputFocus"
                ></textarea>
              </div>
              <button class="emoji-btn" @click="toggleEmojiPicker" type="button" title="插入表情">
                😊
              </button>
            </div>
            <div class="editor-actions">
              <div v-if="replyingTo" class="reply-indicator">
                <span class="reply-label">回复 @{{ getReplyTargetName() }}</span>
                <button class="cancel-reply-btn" @click="cancelReply">
                  <font-awesome-icon icon="times" />
                </button>
              </div>
              <div class="action-buttons">
                <span class="tip">Ctrl + Enter 快速发布</span>
                <span class="char-count">{{ newComment.length }}/300</span>
                <button class="preview-btn" @click="togglePreview">
                  <font-awesome-icon icon="eye" />
                  预览
                </button>
                <button class="submit-btn" :disabled="!newComment.trim()" @click="submitComment">
                  <font-awesome-icon icon="paper-plane" />
                  发布
                </button>
              </div>
            </div>
            <!-- 预览评论 -->
            <div v-if="previewVisible" class="preview-comment">
              <div class="preview-label">预览效果</div>
              <div class="comment-item preview-item">
                <div class="comment-avatar">
                  <div class="avatar-circle">{{ (user.value?.name || user.value?.username) ? (user.value?.name || user.value?.username).charAt(0) : 'U' }}</div>
                </div>
                <div class="comment-body">
                  <div class="comment-header">
                    <span class="comment-author">{{ user.value?.name || user.value?.username || '当前用户' }}</span>
                    <span v-if="replyingTo" class="reply-tag">@{{ getReplyTargetName() }}</span>
                    <span class="comment-time">刚刚</span>
                  </div>
                  <div v-if="newComment.trim()" class="comment-bubble">
                    <div class="comment-content" v-html="renderedPreview"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Emoji 选择器 -->
        <EmojiPicker
          :visible="emojiPickerVisible"
          @select="insertEmoji"
          @close="emojiPickerVisible = false"
        />

        <!-- 评论列表 -->
        <div v-if="Array.isArray(comments) && comments.length > 0" class="comments-list">
          <div v-for="comment in getAllCommentsInOrder()" :key="comment.ID" class="comment-item" :class="{ 'reply-comment': comment.parent_id }">
            <div class="comment-avatar">
              <div class="avatar-circle" :class="{ 'small': comment.parent_id }">{{ comment.username.charAt(0) }}</div>
            </div>
            <div class="comment-body" :class="{ 'selected': replyingTo === comment.ID }">
              <div class="comment-header">
                <span class="comment-author">{{ comment.username }}</span>
                <span v-if="comment.parent_id" class="reply-tag">@{{ getParentCommentUsername(comment.parent_id) }}</span>
                <span class="comment-time">{{ formatCommentTime(comment.CreatedAt) }}</span>
                <button class="hover-reply-btn" @click="startReply(comment.ID, comment.username)">
                  <font-awesome-icon icon="reply" />
                  回复
                </button>
              </div>
              <div class="comment-bubble">
                <!-- 引用原文显示 -->
                <div v-if="comment.quoted_text" class="comment-quoted-text">
                  <div class="quoted-text-label">引用原文：</div>
                  <div
                    :ref="el => setCommentQuotedTextRef(comment.ID, el)"
                    class="quoted-text-content markdown-body"
                    :class="{ 'is-truncated': isCommentQuotedTextTruncated[comment.ID] }"
                    v-html="commentQuotedTexts[comment.ID] || ''"
                  ></div>
                </div>
                <div class="comment-content markdown-body" v-html="formattedCommentContents[comment.ID] || comment.content"></div>
              </div>
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
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick, onActivated } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { getArticleByID } from '@/api/Articles/browse'
import { getCommentsByID } from '@/api/Comments/browse'
import { createComment, deleteComment as deleteCommentAPI } from '@/api/Comments/edit'
import { getLikeStatus, toggleLike } from '@/api/Likes/like'
import { showErrorMessage, showSuccessMessage, showWarningMessage, showCustomMessage } from '@/utils/waifuMessage'
import { generateArticleSEO, updateSEO } from '@/utils/seo'
import RelatedArticles from '@/components/RelatedArticles.vue'
import { estimateReadingTime } from '@/utils/readingTime'
import TextSelectionMenu from '@/components/TextSelectionMenu.vue'
import EmojiPicker from '@/components/EmojiPicker.vue'
import ShareCard from '@/components/ShareCard.vue'
import { renderQuotedText } from '@/utils/renderQuotedText'

// 大型库按需加载，减少首屏 JS 体积
let marked = null
let hljs = null
let protectLatex = null
let restoreAndRenderLatex = null

// 动态加载大型库
async function loadMarkdownLibs () {
  if (marked && hljs && protectLatex && restoreAndRenderLatex) {
    return // 已加载
  }

  try {
    // 并行加载所有库
    const [
      markedModule,
      hljsModule,
      latexModule
    ] = await Promise.all([
      import('marked'),
      import('highlight.js'),
      import('@/utils/latex'),
      import('@/assets/styles/github-highlight.css') // CSS 也需要加载
    ])

    marked = markedModule.marked
    hljs = hljsModule.default
    protectLatex = latexModule.protectLatex
    restoreAndRenderLatex = latexModule.restoreAndRenderLatex
  } catch (error) {
    console.error('加载 Markdown 库失败:', error)
  }
}

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
const router = useRouter()
const user = computed(() => store.state.user)

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

// 精确调整目录位置 - 确保固定在屏幕左侧
const adjustTocPosition = () => {
  nextTick(() => {
    const tocElement = document.querySelector('.floating-toc')
    const contentContainer = document.querySelector('.content-container')

    if (!tocElement || !contentContainer) return

    // 清除任何可能冲突的right属性
    tocElement.style.right = 'auto'

    // 确保目录固定在屏幕左侧
    tocElement.style.left = '20px'
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
    // 排除标识标签，只获取标题的实际文本内容
    // 克隆标题元素，移除标识标签，然后获取文本
    const clonedHeading = heading.cloneNode(true)
    const labelElement = clonedHeading.querySelector('.heading-label')
    if (labelElement) {
      labelElement.remove()
    }
    const text = clonedHeading.textContent.trim()
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

// 点赞功能
const handleLike = async () => {
  // 检查是否已登录
  if (!user.value.isLogged) {
    showCustomMessage('请先登录后再点赞～', 3000)
    return
  }

  // 防止重复点击
  if (likeLoading.value) {
    return
  }

  const articleType = props.type || route.params.type
  const articleID = props.articleId || route.params.id

  if (!articleType || !articleID) {
    showErrorMessage('文章信息不完整')
    return
  }

  likeLoading.value = true

  try {
    const result = await toggleLike(articleType, articleID, user.value, store.state.token)
    isLiked.value = result.isLiked
    likeCount.value = result.likeCount

    // 显示成功提示
    if (result.isLiked) {
      showCustomMessage('点赞成功！感谢你的支持～', 2000)
    } else {
      showCustomMessage('已取消点赞', 2000)
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    // 打印详细的错误信息
    if (error.response) {
      console.error('错误响应:', error.response.data)
      console.error('错误状态:', error.response.status)
      if (error.response.data?.details) {
        console.error('错误详情:', error.response.data.details)
      }
    }
    if (error.response && error.response.status === 401) {
      showCustomMessage('请先登录后再点赞～', 3000)
    } else {
      const errorMsg = error.response?.data?.error || error.response?.data?.details || '点赞操作失败，请稍后重试'
      showErrorMessage(errorMsg)
    }
  } finally {
    likeLoading.value = false
  }
}

// 订阅功能
const handleSubscribe = () => {
  showCustomMessage('RSS功能取消啦，如果需要可以和我说喔', 5000)
}

// 分享功能（复制链接）
const handleShare = async () => {
  const url = window.location.href
  const shareText = `${title.value} - ${url}`
  await copyToClipboard(shareText)
  showCustomMessage('链接已复制到剪贴板，快去分享给朋友吧～', 4000)
}

// 滚动到评论区
const scrollToComments = () => {
  const commentsSection = document.querySelector('.comments-section')
  if (commentsSection) {
    commentsSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
  } else {
    showWarningMessage('评论区域未找到')
  }
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
      const maxScrollTop = tocNav.scrollHeight - tocVisibleHeight

      // 如果目录已经滚动到底部，且当前标题在最后几个位置，不要强制滚动
      const isNearBottom = currentScrollTop >= maxScrollTop - 5 // 允许5px的误差
      const isLastFewItems = closestHeading.index >= toc.value.length - 3

      // 获取目录项的实际位置和尺寸
      const tocItemRect = targetTocLink.getBoundingClientRect()
      const tocNavRect = tocNav.getBoundingClientRect()

      // 计算目录项相对于目录容器的位置
      const itemTop = tocItemRect.top - tocNavRect.top + currentScrollTop
      const itemBottom = itemTop + tocItemRect.height

      // 检查当前标题是否在可视区域内（允许一些误差）
      const tolerance = 2 // 允许2px的误差
      const isItemVisible = itemTop >= currentScrollTop - tolerance &&
                           itemBottom <= currentScrollTop + tocVisibleHeight + tolerance

      // 如果不在可视区域内，则滚动到合适位置
      // 但如果目录已经接近底部且当前是最后几个标题，不要强制滚动
      if (!isItemVisible && !(isNearBottom && isLastFewItems)) {
        let targetScrollTop

        // 如果标题在可视区域上方，滚动到标题位置
        if (itemTop < currentScrollTop) {
          targetScrollTop = itemTop
        } else if (itemBottom > currentScrollTop + tocVisibleHeight) {
          // 如果标题在可视区域下方，滚动到标题底部
          targetScrollTop = itemBottom - tocVisibleHeight
        } else {
          // 如果已经在可视区域内，不需要滚动
          return
        }

        // 处理边界值
        const safeMaxScrollTop = Math.max(0, maxScrollTop)

        // 确保 targetScrollTop 已定义且有效
        if (targetScrollTop !== undefined && !isNaN(targetScrollTop)) {
          const finalScrollTop = Math.max(0, Math.min(targetScrollTop, safeMaxScrollTop))

          // 只有目标位置与当前位置不同时才滚动（允许1px的误差）
          if (Math.abs(finalScrollTop - currentScrollTop) > 1) {
            tocNav.scrollTo({
              top: finalScrollTop,
              behavior: 'instant'
            })
          }
        }
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

const copyToClipboard = async (text = window.location.href) => {
  try {
    await navigator.clipboard.writeText(text)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

// 跳转到编辑页面
const goToEdit = () => {
  const id = props.articleId || route.params.id
  const type = props.type || route.params.type
  router.push({
    path: `/edit/${id}`,
    query: {
      contentType: 'article',
      articleType: type
    }
  })
}

const image = ref('')
const title = ref('')
const tags = ref([])
const time = ref('')
const content = ref('')
const renderedContent = ref('')
const comments = ref([])
const newComment = ref('')
const previewVisible = ref(false)
const replyingTo = ref(null)
const replyTargetName = ref('')
const emojiPickerVisible = ref(false)
const commentTextarea = ref(null)
const viewCount = ref(0)
const toc = ref([]) // 文章目录
const isLiked = ref(false) // 点赞状态
const likeCount = ref(0) // 点赞数
const likeLoading = ref(false) // 点赞操作加载状态

// 文本选择相关
const articleContentRef = ref(null)
const textSelectionMenuVisible = ref(false)
const textSelectionPosition = ref({ x: 0, y: 0 })
const selectedText = ref('')
const isTextHighlighted = ref(false)
const highlightedRange = ref(null)
const shareCardVisible = ref(false)
const shareSelectedText = ref('')
const articleUrl = ref('')
const readingTime = ref(null) // 阅读时间估算
const tocVisible = ref(true) // 目录是否可见
const contentContainer = ref(null) // 内容容器引用
const quotedText = ref('') // 引用的原文
const renderedQuotedText = ref('') // 渲染后的引用文本HTML
const commentQuotedTexts = ref({}) // 存储每个评论的渲染后的引用文本 { commentId: html }
const inputQuotedTextRef = ref(null) // 输入框引用文本的DOM引用
const commentQuotedTextRefs = ref({}) // 评论引用文本的DOM引用 { commentId: element }
const isInputQuotedTextTruncated = ref(false) // 输入框引用文本是否被截断
const isCommentQuotedTextTruncated = ref({}) // 每个评论的引用文本是否被截断 { commentId: boolean }
const fallbackImg = '/images/sunset-mountains.jpg'

// 图片错误回退
const onImgError = (e) => {
  const img = e?.target
  if (img && img.src !== fallbackImg) {
    img.src = fallbackImg
  }
}

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
      // 解析标签（Tags可能是逗号分隔的字符串）
      if (response.Tags && typeof response.Tags === 'string') {
        tags.value = response.Tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0)
      } else if (Array.isArray(response.Tags)) {
        tags.value = response.Tags
      } else {
        tags.value = []
      }
      time.value = response.CreatedAt ? new Date(response.CreatedAt).toLocaleDateString('zh-CN') : '未知时间'
      content.value = response.Content
      viewCount.value = response.viewCount || 0
      likeCount.value = response.likeCount || 0

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
      likeCount.value = response.likeCount || 0

      // 计算阅读时间
      readingTime.value = estimateReadingTime(response.content, response.title)
    }
  } catch (error) {
    console.error('加载文章详情失败:', error)

    // 处理404错误
    if (error.response && error.response.status === 404) {
      // 设置404页面内容（页面内提示“文章不存在”即可，不再让看板娘重复提示）
      image.value = 'https://picsum.photos/id/180/1920/1080'
      title.value = '文章未找到'
      tags.value = []
      time.value = '未知时间'
      content.value = '抱歉，您访问的文章不存在或已被删除。'
      viewCount.value = 0
      readingTime.value = null

      // 不调用 showErrorMessage，避免看板娘弹出“文章不存在”的提示语
      return
    }

    // 处理其他错误
    showErrorMessage('加载文章失败，请稍后重试')
    return
  }

  // 确保库已加载
  await loadMarkdownLibs()
  if (!marked || !protectLatex || !restoreAndRenderLatex) {
    console.error('Markdown 库加载失败')
    showErrorMessage('内容渲染失败，请刷新页面重试')
    return
  }

  // 预处理：保护 LaTeX 公式（在 marked 渲染前处理）
  const { protected: protectedContent, placeholders: latexPlaceholders } = protectLatex(content.value)

  // 预处理：修复 **`code`** 格式（在 marked 渲染前处理）
  // 使用临时标记避免与 marked 的解析冲突
  let processedContent = protectedContent
  const boldCodePlaceholders = new Map()
  let placeholderIndex = 0

  // 匹配 **`code`** 格式，使用临时占位符替换
  processedContent = processedContent.replace(/\*\*`([^`]+)`\*\*/g, (match, codeContent) => {
    const placeholder = `__BOLD_CODE_PLACEHOLDER_${placeholderIndex}__`
    boldCodePlaceholders.set(placeholder, codeContent)
    placeholderIndex++
    return placeholder
  })

  // 渲染 Markdown 内容
  renderedContent.value = marked(processedContent, {
    breaks: false, // 不将换行符转换为 <br>
    gfm: true, // 启用 GitHub Flavored Markdown
    headerIds: false,
    mangle: false,
    pedantic: false, // 禁用严格模式
    sanitize: false, // 不禁用HTML标签
    smartLists: true,
    smartypants: false
  })

  // 后处理：将临时占位符替换为正确的 HTML
  boldCodePlaceholders.forEach((codeContent, placeholder) => {
    renderedContent.value = renderedContent.value.replace(
      new RegExp(placeholder.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g'),
      `<strong><code>${codeContent}</code></strong>`
    )
  })

  // 在博客/随笔文章详情页显示加载成功消息（1/5概率）
  if (props.type === 'blog' || props.type === 'moment') {
    if (window.showMessage && title.value) {
      // 1/5 概率显示消息
      if (Math.random() < 0.2) {
        window.showMessage(`《${title.value}》加载成功，好像没有出现bug，好耶！`, 3000, 10)
      }
    }
  }

  // 后处理：修复没有被正确渲染的粗体语法
  // 修复 **text:** 这样的模式手动转换为 <strong>text:</strong>
  renderedContent.value = renderedContent.value.replace(/\*\*([^*:]+:\**)\*\*/g, '<strong>$1</strong>')

  // 为图片添加内联样式，确保在渲染时就有宽度和高度限制，避免闪烁
  renderedContent.value = renderedContent.value.replace(
    /<img([^>]*)(style="[^"]*")?([^>]*)>/gi,
    (match, before, existingStyle, after) => {
      // 如果已经有style属性，则合并样式
      if (existingStyle) {
        const newStyle = existingStyle.replace(/"/g, '') + '; max-width: 80% !important; max-height: 750px !important; width: auto !important; height: auto !important; display: block !important; margin: 20px auto !important; box-sizing: border-box !important;'
        return `<img${before} style="${newStyle}"${after}>`
      } else {
        // 如果没有style属性，则添加新的style属性
        return `<img${before} style="max-width: 80% !important; max-height: 750px !important; width: auto !important; height: auto !important; display: block !important; margin: 20px auto !important; box-sizing: border-box !important;"${after}>`
      }
    }
  )

  // 恢复并渲染 LaTeX 公式（异步）
  renderedContent.value = await restoreAndRenderLatex(renderedContent.value, latexPlaceholders)

  // 在渲染后通过 nextTick 为正文中的图片添加错误处理，并确保 LaTeX 公式居中
  // 这样即使封面图片和正文图片URL相同，也不会互相影响
  nextTick(() => {
    // 延迟执行，确保 DOM 完全渲染
    setTimeout(() => {
      const markdownBody = document.querySelector('.markdown-body')
      if (markdownBody) {
        const images = markdownBody.querySelectorAll('img')
        images.forEach(img => {
          // 避免重复添加，也避免覆盖封面图片的错误处理
          if (!img.dataset.errorHandlerAdded && !img.classList.contains('article-image')) {
            img.dataset.errorHandlerAdded = 'true'
            img.addEventListener('error', () => {
              img.style.display = 'none'
            })
            // 确保图片样式正确应用（内联样式优先级最高）
            img.style.setProperty('max-width', '80%', 'important')
            img.style.setProperty('max-height', '750px', 'important')
            img.style.setProperty('width', 'auto', 'important')
            img.style.setProperty('height', 'auto', 'important')
            img.style.setProperty('display', 'block', 'important')
            img.style.setProperty('margin', '20px auto', 'important')
            img.style.setProperty('box-sizing', 'border-box', 'important')
          }
        })

        // 确保 LaTeX 公式块居中（在 DOM 更新后再次设置，确保不被覆盖）
        const katexBlocks = markdownBody.querySelectorAll('.katex-block')
        katexBlocks.forEach(block => {
          block.style.setProperty('text-align', 'center', 'important')
          block.style.setProperty('display', 'block', 'important')
          block.style.setProperty('width', '100%', 'important')
          block.style.setProperty('margin', '24px 0', 'important')
          block.style.setProperty('max-width', '100%', 'important')
          // 确保内部的 katex 也居中
          const katex = block.querySelector('.katex')
          if (katex) {
            katex.style.setProperty('text-align', 'center', 'important')
            katex.style.setProperty('display', 'inline-block', 'important')
            katex.style.setProperty('margin', '0 auto', 'important')
          }
          // 确保 katex-block 内部的所有元素都居中
          const allChildren = block.querySelectorAll('*')
          allChildren.forEach(child => {
            child.style.setProperty('text-align', 'center', 'important')
          })
        })

        // 确保包裹 katex-block 的父元素不影响居中
        const parentElements = markdownBody.querySelectorAll('p, div')
        parentElements.forEach(parent => {
          if (parent.querySelector('.katex-block') && !parent.classList.contains('katex-block')) {
            parent.style.setProperty('text-align', 'center', 'important')
          }
        })
      }
    }, 100)
  })

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
  // 直接更新SEO，确保标题及时更新
  updateSEO(seoData)

  // 等待 DOM 更新后手动触发代码高亮和添加复制按钮
  nextTick(() => {
    // 先修复可能残留的 **粗体** 文本节点（不影响已正确渲染的 strong 标签）
    fixResidualBoldInDOM()

    // 延迟执行，确保代码块完全渲染
    setTimeout(async () => {
      // 确保 highlight.js 已加载
      if (!hljs) {
        await loadMarkdownLibs()
      }
      if (hljs) {
        document.querySelectorAll('pre code').forEach((block) => {
          hljs.highlightElement(block)
          addCopyButton(block)
        })
      }

      // 为标题添加标识标签
      addHeadingLabels()

      // 生成目录
      generateTOC()

      // 恢复高亮（在DOM完全渲染后）
      setTimeout(() => {
        restoreHighlights()
      }, 200)
    }, 100) // 延迟100ms确保渲染完成
  })
}

// 加载点赞状态
const loadLikeStatus = async () => {
  try {
    const articleType = props.type || route.params.type
    const articleID = props.articleId || route.params.id

    // 如果ID无效，跳过点赞状态加载
    if (!articleID || articleID === 'undefined' || !articleType) {
      isLiked.value = false
      likeCount.value = 0
      return
    }

    const result = await getLikeStatus(articleType, articleID)
    isLiked.value = result.isLiked || false
    likeCount.value = result.likeCount || 0
  } catch (error) {
    console.error('加载点赞状态失败:', error)
    // 点赞状态加载失败不影响页面显示，静默处理
    isLiked.value = false
    likeCount.value = 0
  }
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
    // 渲染所有评论的引用文本和格式化评论内容
    await nextTick()
    await renderAllCommentQuotedTexts()
    await formatAllCommentContents()
  } catch (error) {
    console.error('加载评论失败:', error)
    comments.value = []
    // 评论加载失败不影响页面显示，静默处理
  }
}

// 开始回复
const startReply = (commentId, username) => {
  replyingTo.value = commentId
  replyTargetName.value = username
  // 聚焦到输入框
  nextTick(() => {
    const textarea = document.querySelector('.comment-editor textarea')
    if (textarea) {
      textarea.focus()
    }
  })
}

// 取消回复
const cancelReply = () => {
  replyingTo.value = null
  replyTargetName.value = ''
}

// 获取回复目标名称
const getReplyTargetName = () => {
  return replyTargetName.value
}

// 处理输入框聚焦
const handleInputFocus = () => {
  // 如果当前没有在回复状态，清空回复状态
  if (!replyingTo.value) {
    cancelReply()
  }
}

// 切换 Emoji 选择器
const toggleEmojiPicker = () => {
  emojiPickerVisible.value = !emojiPickerVisible.value
}

// 插入 Emoji
const insertEmoji = (emoji) => {
  const textarea = commentTextarea.value
  if (!textarea) {
    // 如果没有 textarea ref，直接追加到末尾
    newComment.value += emoji
    return
  }

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = newComment.value

  // 在光标位置插入 Emoji
  newComment.value = text.substring(0, start) + emoji + text.substring(end)

  // 移动光标到插入位置之后
  nextTick(() => {
    textarea.focus()
    const newPosition = start + emoji.length
    textarea.setSelectionRange(newPosition, newPosition)
  })

  // 关闭选择器
  emojiPickerVisible.value = false
}

// 预览切换
const togglePreview = () => {
  previewVisible.value = !previewVisible.value
  if (previewVisible.value) {
    nextTick(async () => {
      // 确保 highlight.js 已加载
      if (!hljs) {
        await loadMarkdownLibs()
      }
      if (hljs) {
        document.querySelectorAll('.preview-comment pre code').forEach((block) => {
          try { hljs.highlightElement(block) } catch (e) {}
        })
      }
    })
  }
}

// 预览渲染（Markdown -> HTML）
// 注意：computed 不能是 async，所以使用 ref + watch 代替
const renderedPreview = ref('')

watch(newComment, async () => {
  if (!newComment.value) {
    renderedPreview.value = ''
    return
  }

  // 先处理换行符，将 \n 转换为 <br>
  const content = newComment.value.replace(/\n/g, '<br>')

  // 确保 marked 已加载
  if (!marked) {
    await loadMarkdownLibs()
  }

  if (!marked) {
    renderedPreview.value = content
    return
  }

  // 然后使用 marked 渲染 Markdown
  renderedPreview.value = marked(content, {
    breaks: false, // 我们已经手动处理了换行
    gfm: true,
    headerIds: false,
    mangle: false,
    pedantic: false,
    sanitize: false,
    smartLists: true,
    smartypants: false
  })
})

// 获取所有评论按正确顺序排列（回复紧跟父评论）
const getAllCommentsInOrder = () => {
  const result = []

  // 遍历所有顶级评论
  for (const comment of comments.value) {
    // 添加顶级评论
    result.push(comment)

    // 添加该评论的所有回复（按时间倒序）
    if (comment.replies && comment.replies.length > 0) {
      // 对回复按时间排序（最新的在前）
      const sortedReplies = [...comment.replies].sort((a, b) =>
        new Date(b.CreatedAt) - new Date(a.CreatedAt)
      )
      result.push(...sortedReplies)
    }
  }

  return result
}

// 获取父评论用户名
const getParentCommentUsername = (parentId) => {
  // 在所有评论中查找父评论（包括顶级评论和回复）
  const allComments = getAllCommentsInOrder()
  for (const comment of allComments) {
    if (comment.ID === parentId) {
      return comment.username
    }
  }
  return '未知用户'
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) return

  const id = props.articleId || route.params.id

  if (!id || id === 'undefined') {
    showErrorMessage('文章ID无效，无法提交评论')
    return
  }

  if (!user.value || !user.value.isLogged || !store.state.token) {
    console.log('登录状态检查失败:', {
      user: user.value,
      isLogged: user.value?.isLogged,
      token: store.state.token,
      tokenLength: store.state.token?.length
    })
    showErrorMessage('401')
    return
  }

  try {
    // 如果是回复，使用parentId
    const parentId = replyingTo.value || null
    // 传递引用文本
    const quoted = quotedText.value || null
    await createComment(user.value, id, props.type, newComment.value, parentId, store.state.token, quoted)

    newComment.value = ''
    cancelReply() // 清空回复状态
    clearQuotedText() // 清空引用文本
    await loadComments()
    // 重新渲染引用文本和格式化评论内容（包括新提交的评论）
    await renderAllCommentQuotedTexts()
    await formatAllCommentContents()
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
    await deleteCommentAPI(user, commentId, store.state.token)
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
// 清理选中文本中的UI元素（代码块复制按钮、标题标签等）
const cleanSelectedText = (range) => {
  if (!range) return ''

  try {
    // 克隆range的内容，避免修改原始DOM
    const clonedContents = range.cloneContents()

    // 创建一个临时容器来操作克隆的内容
    const tempDiv = document.createElement('div')
    tempDiv.appendChild(clonedContents)

    // 移除代码块复制按钮（包括emoji图标）
    const copyButtons = tempDiv.querySelectorAll('.copy-btn')
    copyButtons.forEach(btn => btn.remove())

    // 移除标题标签（H1, H2等）
    const headingLabels = tempDiv.querySelectorAll('.heading-label')
    headingLabels.forEach(label => label.remove())

    // 移除其他可能的UI元素
    const uiElements = tempDiv.querySelectorAll('[class*="toolbox"], [class*="button"], [class*="btn"]')
    uiElements.forEach(el => {
      // 只移除明显的UI元素，保留内容相关的元素
      if (el.classList.contains('copy-btn') ||
          el.classList.contains('heading-label') ||
          el.classList.contains('code-toolbox') ||
          el.classList.contains('expand-button') ||
          el.classList.contains('copy-button')) {
        el.remove()
      }
    })

    // 尝试从原始Markdown中提取对应的文本（包括LaTeX公式和代码块）
    // 如果选中的内容包含LaTeX公式或代码块的渲染结果，尝试从原始内容中匹配
    const katexElements = tempDiv.querySelectorAll('.katex, .katex-block')
    const codeBlocks = tempDiv.querySelectorAll('pre code, code')
    let extractedText = ''

    // 优先处理代码块（因为代码块格式更明显，匹配更准确）
    if (codeBlocks.length > 0 && content.value) {
      // 如果包含代码块，尝试从原始Markdown中提取
      const codeText = Array.from(codeBlocks).map(block => block.textContent).join('\n')
      const lines = content.value.split('\n')
      let foundMatch = false

      // 查找包含代码块的行（```标记）
      for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim()
        if (line.startsWith('```')) {
          // 找到代码块开始，提取整个代码块
          let codeBlock = line + '\n'

          for (let j = i + 1; j < lines.length; j++) {
            const currentLine = lines[j]
            if (currentLine.trim() === '```') {
              codeBlock += currentLine
              extractedText = codeBlock
              foundMatch = true
              break
            } else {
              codeBlock += currentLine + '\n'
            }
          }

          if (foundMatch) {
            // 验证提取的代码块内容是否匹配
            const extractedCode = codeBlock.replace(/```[\w]*\n?/g, '').replace(/```/g, '').trim()
            if (extractedCode.includes(codeText.substring(0, Math.min(50, codeText.length)))) {
              return extractedText
            }
          }
        }
      }
    }

    // 处理LaTeX公式
    if (katexElements.length > 0 && content.value && !extractedText) {
      // 如果包含LaTeX公式，尝试从原始Markdown中提取
      // 获取选中文本的纯文本版本（用于匹配）
      const plainText = tempDiv.textContent || tempDiv.innerText || ''
      const cleanPlainText = plainText.replace(/\s+/g, ' ').trim()

      // 尝试在原始Markdown中找到包含LaTeX公式的段落
      // 策略：查找包含LaTeX公式标记（$$或$）的行，并尝试匹配上下文
      const lines = content.value.split('\n')
      let foundMatch = false
      const searchText = cleanPlainText.length > 30 ? cleanPlainText.substring(0, 30) : cleanPlainText

      // 首先尝试找到包含LaTeX公式的行
      for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim()
        // 如果这一行包含LaTeX公式标记
        if (line.includes('$$') || (line.includes('$') && !line.startsWith('$'))) {
          // 检查这一行或相邻行是否包含选中的文本片段
          const contextLines = []
          // 收集当前行和前后各2行作为上下文
          for (let j = Math.max(0, i - 2); j <= Math.min(lines.length - 1, i + 2); j++) {
            contextLines.push(lines[j])
          }
          const context = contextLines.join('\n')

          // 如果上下文中包含选中的文本片段，提取包含LaTeX公式的行
          if (context.includes(searchText) || cleanPlainText.length < 20) {
            // 提取包含LaTeX公式的完整段落
            // 对于块级公式（$$），提取完整的公式块
            if (line.includes('$$')) {
              // 查找公式块的开始和结束
              let formulaBlock = ''
              let inBlock = false
              for (let j = i; j < lines.length; j++) {
                const currentLine = lines[j]
                if (currentLine.includes('$$')) {
                  formulaBlock += currentLine + '\n'
                  // 如果遇到第二个$$，结束块级公式
                  if ((currentLine.match(/\$\$/g) || []).length >= 2) {
                    break
                  } else if (!inBlock) {
                    inBlock = true
                  } else {
                    break
                  }
                } else if (inBlock) {
                  formulaBlock += currentLine + '\n'
                } else {
                  break
                }
              }
              extractedText = formulaBlock.trim()
            } else {
              // 行内公式，提取整行
              extractedText = line
            }
            foundMatch = true
            break
          }
        }
      }

      // 如果找到了匹配，使用提取的文本；否则使用清理后的HTML文本
      if (foundMatch && extractedText) {
        return extractedText
      }
    }

    // 提取清理后的文本，保留换行
    // 使用 innerHTML 获取 HTML 结构，然后转换为 Markdown 格式以保留换行
    let cleanedText = ''

    // 尝试保留换行：将块级元素的换行转换为换行符
    const blockElements = tempDiv.querySelectorAll('p, div, br, pre, li, h1, h2, h3, h4, h5, h6')
    if (blockElements.length > 0) {
      // 如果有块级元素，尝试保留结构
      // 将 <br> 转换为换行
      const htmlContent = tempDiv.innerHTML
      cleanedText = htmlContent
        .replace(/<br\s*\/?>/gi, '\n') // <br> 转换为换行
        .replace(/<\/p>/gi, '\n') // </p> 后添加换行
        .replace(/<\/div>/gi, '\n') // </div> 后添加换行
        .replace(/<\/li>/gi, '\n') // </li> 后添加换行
        .replace(/<\/h[1-6]>/gi, '\n') // 标题后添加换行
        .replace(/<[^>]+>/g, '') // 移除所有HTML标签
        .replace(/\n{3,}/g, '\n\n') // 多个连续换行替换为两个换行
        .trim()
    } else {
      // 如果没有块级元素，使用 textContent 但保留换行
      cleanedText = tempDiv.textContent || tempDiv.innerText || ''
      // 保留换行，只清理多余的空白字符
      cleanedText = cleanedText
        .replace(/[ \t]+/g, ' ') // 将多个连续空格/制表符替换为单个空格
        .replace(/\n{3,}/g, '\n\n') // 多个连续换行替换为两个换行
        .trim()
    }

    return cleanedText
  } catch (error) {
    console.error('清理选中文本失败:', error)
    // 如果清理失败，返回原始文本（通过toString获取）
    return range.toString().trim()
  }
}

// 文本选择处理函数
const handleTextSelection = (e) => {
  // 如果点击的是文本选择菜单，不处理文本选择
  if (e.target && (e.target.closest('.text-selection-menu') || e.target.closest('.menu-btn'))) {
    return
  }

  const selection = window.getSelection()

  // 如果选择为空或不在文章内容区域内，隐藏菜单
  if (!selection.rangeCount || !articleContentRef.value) {
    textSelectionMenuVisible.value = false
    return
  }

  // 检查选择是否在 markdown-body 内
  const markdownBody = articleContentRef.value.querySelector('.markdown-body')
  if (!markdownBody || !markdownBody.contains(selection.anchorNode)) {
    textSelectionMenuVisible.value = false
    return
  }

  // 获取选择的位置
  const range = selection.getRangeAt(0).cloneRange() // 克隆 range，避免选择被清除时丢失

  // 清理选中文本，移除UI元素，并尝试提取原始Markdown（包括LaTeX公式）
  const selectedTextValue = cleanSelectedText(range)

  if (!selectedTextValue) {
    textSelectionMenuVisible.value = false
    return
  }

  const rect = range.getBoundingClientRect()

  // 设置菜单位置（鼠标位置）
  textSelectionPosition.value = {
    x: rect.left + rect.width / 2,
    y: rect.top
  }

  selectedText.value = selectedTextValue
  highlightedRange.value = range

  // 检查是否已高亮
  isTextHighlighted.value = checkIfHighlighted(range)

  textSelectionMenuVisible.value = true
}

// 清除选择（只隐藏菜单，不清除高亮）
// 注意：此函数目前未被直接使用，但保留以备将来需要
// const clearTextSelection = () => {
//   textSelectionMenuVisible.value = false
//   selectedText.value = ''
//   highlightedRange.value = null
// }

// 复制文本
const handleTextCopy = async (text) => {
  try {
    await copyToClipboard(text)
    showSuccessMessage('copy')
    // 只隐藏菜单，不清除高亮
    textSelectionMenuVisible.value = false
    selectedText.value = ''
    highlightedRange.value = null
  } catch (error) {
    showErrorMessage('复制失败')
  }
}

// 检查选择是否已高亮
const checkIfHighlighted = (range) => {
  if (!range) return false
  const container = range.commonAncestorContainer
  const parent = container.nodeType === Node.TEXT_NODE ? container.parentElement : container
  return parent?.closest('.text-highlight') !== null
}

// 高亮文本 - 使用 span 标签，但设置 display: contents 以避免破坏布局
// 获取范围内的所有文本节点
const getTextNodesInRange = (range) => {
  const textNodes = []
  const walker = document.createTreeWalker(
    range.commonAncestorContainer.nodeType === Node.TEXT_NODE
      ? range.commonAncestorContainer.parentNode
      : range.commonAncestorContainer,
    NodeFilter.SHOW_TEXT,
    null
  )

  let node
  while ((node = walker.nextNode())) {
    if (range.intersectsNode(node)) {
      textNodes.push(node)
    }
  }
  return textNodes
}

// 检查文本节点是否已经高亮（或其祖先元素是否高亮）
const isTextNodeHighlighted = (textNode) => {
  let current = textNode.parentElement || textNode.parentNode
  while (current) {
    // 跳过 highlight.js 添加的元素（如 hljs 相关的类）
    if (current.classList) {
      if (current.classList.contains('text-highlight')) {
        return true
      }
      // 跳过 highlight.js 的 span（通常有 hljs-* 类）
      if (Array.from(current.classList).some(cls => cls.startsWith('hljs-'))) {
        // 继续向上查找，但不把 hljs span 当作高亮
        current = current.parentElement || current.parentNode
        continue
      }
    }
    // 如果到达了文章内容容器的边界，停止搜索
    if (current === articleContentRef.value ||
        (current.nodeType === Node.ELEMENT_NODE && current.classList?.contains('markdown-body'))) {
      break
    }
    current = current.parentElement || current.parentNode
  }
  return false
}

// 获取范围内的所有高亮span元素
const getHighlightSpansInRange = (range) => {
  const highlights = []
  const walker = document.createTreeWalker(
    range.commonAncestorContainer.nodeType === Node.TEXT_NODE
      ? range.commonAncestorContainer.parentNode
      : range.commonAncestorContainer,
    NodeFilter.SHOW_ELEMENT,
    {
      acceptNode: (node) => {
        if (node.classList && node.classList.contains('text-highlight')) {
          return NodeFilter.FILTER_ACCEPT
        }
        return NodeFilter.FILTER_SKIP
      }
    }
  )

  let node
  while ((node = walker.nextNode())) {
    if (range.intersectsNode(node)) {
      highlights.push(node)
    }
  }
  return highlights
}

// 检查选择范围是否完全在高亮区域内
const isRangeFullyHighlighted = (range) => {
  // 获取范围内的所有高亮span
  const highlights = getHighlightSpansInRange(range)
  if (highlights.length === 0) return false

  // 创建一个范围来检查所有文本节点
  const textNodes = getTextNodesInRange(range)
  if (textNodes.length === 0) return false

  // 检查每个文本节点是否都在某个高亮span内
  for (const textNode of textNodes) {
    let startOffset = 0
    let endOffset = textNode.textContent.length

    if (textNode === range.startContainer) {
      startOffset = range.startOffset
    }
    if (textNode === range.endContainer) {
      endOffset = range.endOffset
    }

    // 检查这个文本节点（或部分）是否在高亮span内
    let isCovered = false
    for (const highlight of highlights) {
      const highlightRange = document.createRange()
      highlightRange.selectNodeContents(highlight)

      // 检查文本节点是否与高亮范围相交
      if (highlightRange.intersectsNode(textNode)) {
        // 进一步检查偏移量是否在范围内
        const nodeRange = document.createRange()
        nodeRange.setStart(textNode, startOffset)
        nodeRange.setEnd(textNode, endOffset)

        // 检查这个范围是否完全在高亮范围内
        if (highlightRange.compareBoundaryPoints(Range.START_TO_START, nodeRange) <= 0 &&
            highlightRange.compareBoundaryPoints(Range.END_TO_END, nodeRange) >= 0) {
          isCovered = true
          break
        }
      }
    }

    if (!isCovered) {
      return false
    }
  }

  return true
}

// 高亮文本 - 支持智能合并，避免嵌套
const handleTextHighlight = async (text) => {
  if (!highlightedRange.value) return

  // 使用当前选择重新获取 range（更可靠）
  const selection = window.getSelection()
  let range = null

  if (selection.rangeCount > 0) {
    range = selection.getRangeAt(0).cloneRange()
  } else if (highlightedRange.value) {
    range = highlightedRange.value.cloneRange()
  }

  if (!range || range.collapsed) {
    showCustomMessage('请先选择文本')
    return
  }

  // 检查是否完全在高亮区域内
  if (isRangeFullyHighlighted(range)) {
    // 完全在高亮内，不进行任何操作，只隐藏菜单
    textSelectionMenuVisible.value = false
    selectedText.value = ''
    highlightedRange.value = null
    window.getSelection().removeAllRanges()
    return
  }

  // 不移除重叠的高亮，而是只高亮未高亮的部分
  // 已高亮的部分会通过检查逻辑自动跳过

  // 检查是否包含块级元素（列表、标题等）
  const blockElements = ['P', 'DIV', 'LI', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6', 'UL', 'OL', 'DL', 'DT', 'DD', 'BLOCKQUOTE', 'PRE', 'CODE', 'TABLE', 'TR', 'TD', 'TH', 'THEAD', 'TBODY']
  const containsBlockElement = () => {
    const commonAncestor = range.commonAncestorContainer
    const startContainer = range.startContainer
    const endContainer = range.endContainer

    const checkNode = (node) => {
      let current = node.nodeType === Node.TEXT_NODE ? node.parentElement : node
      while (current && current !== commonAncestor) {
        if (blockElements.includes(current.tagName)) {
          return true
        }
        current = current.parentElement
      }
      return false
    }

    return checkNode(startContainer) || checkNode(endContainer)
  }

  // 如果包含块级元素，逐行处理
  if (containsBlockElement()) {
    const textNodes = getTextNodesInRange(range)
    let hasHighlight = false

    textNodes.forEach(textNode => {
      // 计算与原始range的交集
      let startOffset = 0
      let endOffset = textNode.textContent.length

      if (textNode === range.startContainer) {
        startOffset = range.startOffset
      }
      if (textNode === range.endContainer) {
        endOffset = range.endOffset
      }

      if (endOffset > startOffset) {
        // 检查文本节点是否在高亮span内，以及高亮span覆盖的范围
        // 注意：跳过 highlight.js 添加的元素
        let current = textNode.parentElement || textNode.parentNode
        let highlightSpan = null
        while (current) {
          if (current.classList && current.classList.contains('text-highlight')) {
            highlightSpan = current
            break
          }
          // 跳过 highlight.js 的 span（通常有 hljs-* 类）
          if (current.classList && Array.from(current.classList).some(cls => cls.startsWith('hljs-'))) {
            current = current.parentElement || current.parentNode
            continue
          }
          if (current === articleContentRef.value ||
              (current.nodeType === Node.ELEMENT_NODE && current.classList?.contains('markdown-body'))) {
            break
          }
          current = current.parentElement || current.parentNode
        }

        let highlightStart = null
        let highlightEnd = null

        if (highlightSpan) {
          // 获取高亮span的范围
          const highlightRange = document.createRange()
          highlightRange.selectNodeContents(highlightSpan)

          // 检查文本节点是否在高亮span内
          if (highlightRange.intersectsNode(textNode)) {
            // 创建文本节点的范围用于比较
            const textNodeRange = document.createRange()
            textNodeRange.selectNodeContents(textNode)

            // 检查高亮span的起始和结束位置
            const startContainer = highlightRange.startContainer
            const startOffset = highlightRange.startOffset
            const endContainer = highlightRange.endContainer
            const endOffset = highlightRange.endOffset

            // 如果文本节点是高亮span的起始容器
            if (startContainer === textNode) {
              highlightStart = startOffset
            } else {
              // 检查高亮span是否从文本节点之前或开始位置开始
              const startCompare = highlightRange.compareBoundaryPoints(Range.START_TO_START, textNodeRange)
              if (startCompare <= 0) {
                // 高亮span从文本节点之前或开始位置开始
                highlightStart = 0
              }
            }

            // 如果文本节点是高亮span的结束容器
            if (endContainer === textNode) {
              highlightEnd = endOffset
            } else {
              // 检查高亮span是否延伸到文本节点之后或结束位置
              const endCompare = highlightRange.compareBoundaryPoints(Range.END_TO_END, textNodeRange)
              if (endCompare >= 0) {
                // 高亮span延伸到文本节点之后或结束位置
                highlightEnd = textNode.textContent.length
              }
            }

            // 如果还没有确定，检查文本节点是否完全在高亮span内
            if (highlightStart === null || highlightEnd === null) {
              // 检查文本节点是否完全被高亮span包含
              const startCompare = highlightRange.compareBoundaryPoints(Range.START_TO_START, textNodeRange)
              const endCompare = highlightRange.compareBoundaryPoints(Range.END_TO_END, textNodeRange)

              if (startCompare <= 0 && endCompare >= 0) {
                // 文本节点完全在高亮span内
                highlightStart = 0
                highlightEnd = textNode.textContent.length
              }
            }
          }
        }

        // 计算需要高亮的范围（排除已高亮的部分）
        let needHighlightStart = startOffset
        let needHighlightEnd = endOffset

        if (highlightStart !== null && highlightEnd !== null) {
          // 如果选择范围与已高亮范围有重叠
          if (startOffset < highlightEnd && endOffset > highlightStart) {
            // 如果选择范围完全在已高亮范围内，跳过
            if (startOffset >= highlightStart && endOffset <= highlightEnd) {
              return
            }

            // 如果选择范围部分重叠，只高亮未重叠的部分
            if (startOffset < highlightStart && endOffset > highlightStart) {
              // 选择范围从已高亮之前开始，高亮前面的部分
              needHighlightEnd = Math.min(endOffset, highlightStart)
            } else if (startOffset < highlightEnd && endOffset > highlightEnd) {
              // 选择范围延伸到已高亮之后，高亮后面的部分
              needHighlightStart = Math.max(startOffset, highlightEnd)
            } else if (startOffset >= highlightStart && endOffset <= highlightEnd) {
              // 完全在已高亮内，跳过
              return
            }
          }
        }

        // 只高亮未高亮的部分
        if (needHighlightEnd > needHighlightStart) {
          try {
            const highlightRange = document.createRange()
            highlightRange.setStart(textNode, needHighlightStart)
            highlightRange.setEnd(textNode, needHighlightEnd)

            // 再次检查：确保要高亮的范围不在已高亮的span内
            // 跳过 highlight.js 添加的元素
            let parent = textNode.parentElement || textNode.parentNode
            while (parent) {
              if (parent.classList) {
                if (parent.classList.contains('text-highlight')) {
                  // 要高亮的文本已经在高亮span内，跳过
                  return
                }
                // 跳过 highlight.js 的 span（通常有 hljs-* 类）
                if (Array.from(parent.classList).some(cls => cls.startsWith('hljs-'))) {
                  parent = parent.parentElement || parent.parentNode
                  continue
                }
              }
              if (parent === articleContentRef.value ||
                  (parent.nodeType === Node.ELEMENT_NODE && parent.classList?.contains('markdown-body'))) {
                break
              }
              parent = parent.parentElement || parent.parentNode
            }

            const span = document.createElement('span')
            span.className = 'text-highlight'
            span.style.setProperty('display', 'inline', 'important')
            span.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
            span.style.setProperty('padding', '0', 'important')
            span.style.setProperty('margin', '0', 'important')
            span.style.setProperty('line-height', 'inherit', 'important')
            span.style.setProperty('vertical-align', 'baseline', 'important')

            highlightRange.surroundContents(span)
            hasHighlight = true
          } catch (e) {
            // 如果 surroundContents 失败，手动分割文本节点
            try {
              const text = textNode.textContent
              const beforeText = text.substring(0, needHighlightStart)
              const highlightText = text.substring(needHighlightStart, needHighlightEnd)
              const afterText = text.substring(needHighlightEnd)
              const parent = textNode.parentNode

              if (parent && highlightText.trim()) {
                if (beforeText) {
                  parent.insertBefore(document.createTextNode(beforeText), textNode)
                }
                const span = document.createElement('span')
                span.className = 'text-highlight'
                span.style.setProperty('display', 'inline', 'important')
                span.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
                span.style.setProperty('padding', '0', 'important')
                span.style.setProperty('margin', '0', 'important')
                span.style.setProperty('line-height', 'inherit', 'important')
                span.style.setProperty('vertical-align', 'baseline', 'important')
                span.textContent = highlightText
                parent.insertBefore(span, textNode)
                if (afterText) {
                  parent.insertBefore(document.createTextNode(afterText), textNode)
                }
                parent.removeChild(textNode)
                hasHighlight = true
              }
            } catch (err) {
              console.error('高亮文本节点失败:', err)
            }
          }
        }
      }
    })

    if (hasHighlight) {
      saveHighlights()
      window.getSelection().removeAllRanges()
      textSelectionMenuVisible.value = false
      selectedText.value = ''
      highlightedRange.value = null
    }
  } else {
    // 普通文本，检查是否已高亮，只高亮未高亮的部分
    // 检查选择范围是否完全在高亮内
    const textNodes = getTextNodesInRange(range)
    let allHighlighted = true
    for (const textNode of textNodes) {
      if (!isTextNodeHighlighted(textNode)) {
        allHighlighted = false
        break
      }
    }

    if (allHighlighted) {
      // 全部已高亮，不做操作
      textSelectionMenuVisible.value = false
      selectedText.value = ''
      highlightedRange.value = null
      window.getSelection().removeAllRanges()
      return
    }

    // 检查是否在已高亮的span内（避免嵌套高亮导致颜色变深）
    const commonAncestor = range.commonAncestorContainer
    let parent = commonAncestor.nodeType === Node.TEXT_NODE ? commonAncestor.parentElement : commonAncestor
    while (parent && parent !== articleContentRef.value) {
      if (parent.classList && parent.classList.contains('text-highlight')) {
        // 选择范围在已高亮的span内，不做操作（避免嵌套）
        textSelectionMenuVisible.value = false
        selectedText.value = ''
        highlightedRange.value = null
        window.getSelection().removeAllRanges()
        return
      }
      parent = parent.parentElement
    }

    try {
      const span = document.createElement('span')
      span.className = 'text-highlight'
      span.style.setProperty('display', 'inline', 'important')
      span.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
      span.style.setProperty('padding', '0', 'important')
      span.style.setProperty('margin', '0', 'important')
      span.style.setProperty('line-height', 'inherit', 'important')
      span.style.setProperty('vertical-align', 'baseline', 'important')
      range.surroundContents(span)

      if (span.parentNode) {
        saveHighlights()
        window.getSelection().removeAllRanges()
        textSelectionMenuVisible.value = false
        selectedText.value = ''
        highlightedRange.value = null
      }
    } catch (error) {
      // 如果 surroundContents 失败（跨节点选择），使用 extractContents + insertBefore
      try {
        // 在 extractContents 之前保存插入位置信息
        const startContainer = range.startContainer
        const startOffset = range.startOffset
        const commonAncestor = range.commonAncestorContainer

        // 保存父节点引用（在 extractContents 之前）
        let insertParent = null
        let insertBeforeNode = null

        if (startContainer.nodeType === Node.TEXT_NODE) {
          insertParent = startContainer.parentNode
          // 保存下一个兄弟节点作为插入参考（extractContents 后可能仍然存在）
          insertBeforeNode = startContainer.nextSibling
        } else {
          insertParent = startContainer
          if (startOffset < startContainer.childNodes.length) {
            insertBeforeNode = startContainer.childNodes[startOffset]
          }
        }

        // 提取内容
        const contents = range.extractContents()
        if (contents) {
          const span = document.createElement('span')
          span.className = 'text-highlight'
          span.style.setProperty('display', 'inline', 'important')
          span.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
          span.style.setProperty('padding', '0', 'important')
          span.style.setProperty('margin', '0', 'important')
          span.style.setProperty('line-height', 'inherit', 'important')
          span.style.setProperty('vertical-align', 'baseline', 'important')
          span.appendChild(contents)

          // 尝试插入节点
          try {
            // 首先尝试使用保存的父节点信息
            if (insertParent) {
              insertParent.insertBefore(span, insertBeforeNode)
            } else {
              // 如果父节点不存在，尝试使用 range 的当前状态
              try {
                range.insertNode(span)
              } catch (rangeError) {
                // 最后的尝试：使用 commonAncestor
                if (commonAncestor && commonAncestor.nodeType === Node.ELEMENT_NODE) {
                  if (commonAncestor.firstChild) {
                    commonAncestor.insertBefore(span, commonAncestor.firstChild)
                  } else {
                    commonAncestor.appendChild(span)
                  }
                } else {
                  throw new Error('无法找到有效的插入位置')
                }
              }
            }
          } catch (insertError) {
            // 如果所有方法都失败，尝试在 commonAncestor 中查找
            if (commonAncestor && commonAncestor.nodeType === Node.ELEMENT_NODE) {
              // 查找 commonAncestor 中的第一个文本节点或元素节点作为插入参考
              const walker = document.createTreeWalker(
                commonAncestor,
                NodeFilter.SHOW_TEXT | NodeFilter.SHOW_ELEMENT,
                null
              )
              const firstNode = walker.nextNode()
              if (firstNode) {
                if (firstNode.nodeType === Node.TEXT_NODE) {
                  firstNode.parentNode?.insertBefore(span, firstNode)
                } else {
                  firstNode.insertBefore(span, firstNode.firstChild)
                }
              } else {
                commonAncestor.appendChild(span)
              }
            } else {
              throw new Error('无法找到有效的插入位置')
            }
          }

          if (span.parentNode) {
            saveHighlights()
            window.getSelection().removeAllRanges()
            textSelectionMenuVisible.value = false
            selectedText.value = ''
            highlightedRange.value = null
          }
        }
      } catch (extractError) {
        console.error('高亮失败:', extractError)
      }
    }
  }
}

// 移除高亮（仅在切换文章时使用）
const removeHighlight = () => {
  const highlights = articleContentRef.value?.querySelectorAll('.text-highlight')
  if (highlights) {
    highlights.forEach(highlight => {
      const parent = highlight.parentNode
      if (parent) {
        // 将span的内容提取出来，替换span本身
        while (highlight.firstChild) {
          parent.insertBefore(highlight.firstChild, highlight)
        }
        parent.removeChild(highlight)
        parent.normalize() // 合并相邻的文本节点
      }
    })
  }
  saveHighlights() // 保存高亮状态（清除后保存）
}

// 保存高亮到 sessionStorage（仅在同一会话中保持）
const saveHighlights = () => {
  if (!articleContentRef.value) return

  const id = props.articleId || route.params.id
  if (!id || id === 'undefined') return

  const markdownBody = articleContentRef.value.querySelector('.markdown-body')
  if (!markdownBody) return

  // 获取所有高亮元素及其文本内容和位置信息
  const highlights = markdownBody.querySelectorAll('.text-highlight')
  const highlightData = Array.from(highlights).map(highlight => highlight.textContent.trim()).filter(text => text.length > 0)

  // 保存到 sessionStorage，key 为文章ID
  const storageKey = `article-highlights-${props.type}-${id}`
  if (highlightData.length > 0) {
    sessionStorage.setItem(storageKey, JSON.stringify(highlightData))
  } else {
    sessionStorage.removeItem(storageKey)
  }
}

// 恢复高亮（从 sessionStorage 恢复）
const restoreHighlights = () => {
  if (!articleContentRef.value) return

  const id = props.articleId || route.params.id
  if (!id || id === 'undefined') return

  const markdownBody = articleContentRef.value.querySelector('.markdown-body')
  if (!markdownBody) return

  // 从 sessionStorage 读取高亮数据
  const storageKey = `article-highlights-${props.type}-${id}`
  const savedHighlights = sessionStorage.getItem(storageKey)
  if (!savedHighlights) return

  try {
    const highlightTexts = JSON.parse(savedHighlights)
    if (!Array.isArray(highlightTexts) || highlightTexts.length === 0) return

    // 获取完整的文本内容（用于匹配）
    const fullText = markdownBody.textContent || markdownBody.innerText

    // 对每个高亮文本，在完整文本中查找并恢复
    highlightTexts.forEach(highlightText => {
      if (!highlightText || highlightText.trim() === '') return

      // 检查文本是否存在于当前DOM中
      if (fullText.indexOf(highlightText) === -1) {
        return // 文本不存在，跳过
      }

      // 创建 Range 对象来查找和标记文本
      const range = document.createRange()

      // 遍历所有文本节点，查找匹配的文本
      const walker = document.createTreeWalker(
        markdownBody,
        NodeFilter.SHOW_TEXT,
        null
      )

      let found = false
      let node
      while ((node = walker.nextNode()) && !found) {
        const text = node.textContent
        const index = text.indexOf(highlightText)

        if (index !== -1) {
          // 检查是否已经高亮
          const parent = node.parentElement
          if (parent && parent.classList.contains('text-highlight')) {
            found = true
            continue // 已经高亮，跳过
          }

          // 创建高亮 span
          try {
            range.setStart(node, index)
            range.setEnd(node, index + highlightText.length)

            const span = document.createElement('span')
            span.className = 'text-highlight'
            span.style.setProperty('display', 'inline', 'important')
            span.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
            span.style.setProperty('padding', '0', 'important')
            span.style.setProperty('margin', '0', 'important')
            span.style.setProperty('line-height', 'inherit', 'important')
            span.style.setProperty('vertical-align', 'baseline', 'important')

            try {
              range.surroundContents(span)
            } catch (error) {
              // 如果 surroundContents 失败，使用另一种方法
              const contents = range.extractContents()
              span.appendChild(contents)
              range.insertNode(span)
            }

            found = true
            break
          } catch (error) {
            console.error('恢复高亮失败:', error)
          }
        }
      }
    })

    // 恢复后，确保所有高亮元素都有正确的样式
    setTimeout(() => {
      const highlights = markdownBody.querySelectorAll('.text-highlight')
      highlights.forEach(hl => {
        hl.style.setProperty('display', 'inline', 'important')
        hl.style.setProperty('background-color', 'rgba(168, 85, 247, 0.4)', 'important')
        hl.style.setProperty('padding', '0', 'important')
        hl.style.setProperty('margin', '0', 'important')
        hl.style.setProperty('line-height', 'inherit', 'important')
        hl.style.setProperty('vertical-align', 'baseline', 'important')
      })
      console.log('恢复高亮数量:', highlights.length)
    }, 100)
  } catch (error) {
    console.error('恢复高亮失败:', error)
  }
}

// 分享文本
// 处理文本评论
const handleTextComment = async (text) => {
  if (!text || !text.trim()) return

  // 保存引用的原文
  quotedText.value = text.trim()

  // 检查是否包含代码块或LaTeX公式
  const hasCodeBlock = quotedText.value.includes('```')
  const hasLatex = quotedText.value.includes('$$') || quotedText.value.includes('$')

  // 如果包含代码块或LaTeX公式，使用renderQuotedText渲染
  if (hasCodeBlock || hasLatex) {
    renderedQuotedText.value = await renderQuotedText(quotedText.value, 250)
  } else {
    // 否则只保留换行，不渲染Markdown格式
    // 将换行符转换为 <br>，并转义HTML特殊字符
    const escapedText = quotedText.value
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;')
    renderedQuotedText.value = escapedText.replace(/\n/g, '<br>')
  }

  // 隐藏文本选择菜单
  textSelectionMenuVisible.value = false
  selectedText.value = ''
  highlightedRange.value = null

  // 跳转到评论区
  await nextTick()
  scrollToComments()

  // 检测输入框引用文本是否溢出
  await nextTick()
  checkInputQuotedTextOverflow()

  // 聚焦到评论输入框
  setTimeout(() => {
    const textarea = document.querySelector('.comment-input')
    if (textarea) {
      textarea.focus()
    }
  }, 300)
}

// 检测输入框引用文本是否溢出
const checkInputQuotedTextOverflow = () => {
  if (!inputQuotedTextRef.value) {
    isInputQuotedTextTruncated.value = false
    return
  }
  const element = inputQuotedTextRef.value
  isInputQuotedTextTruncated.value = element.scrollHeight > element.clientHeight || element.scrollWidth > element.clientWidth
}

// 设置评论引用文本的DOM引用
const setCommentQuotedTextRef = (commentId, el) => {
  if (el) {
    commentQuotedTextRefs.value[commentId] = el
    // 检测是否溢出
    nextTick(() => {
      checkCommentQuotedTextOverflow(commentId)
    })
  }
}

// 检测评论引用文本是否溢出
const checkCommentQuotedTextOverflow = (commentId) => {
  const element = commentQuotedTextRefs.value[commentId]
  if (!element) {
    isCommentQuotedTextTruncated.value[commentId] = false
    return
  }
  isCommentQuotedTextTruncated.value[commentId] = element.scrollHeight > element.clientHeight || element.scrollWidth > element.clientWidth
}

// 清除引用文本
const clearQuotedText = () => {
  quotedText.value = ''
  renderedQuotedText.value = ''
}

// 渲染所有评论的引用文本（如果包含代码块或LaTeX公式则渲染，否则只保留换行）
const renderAllCommentQuotedTexts = async () => {
  const allComments = getAllCommentsInOrder()
  const renderPromises = allComments
    .filter(comment => comment.quoted_text && comment.quoted_text.trim())
    .map(async (comment) => {
      if (!commentQuotedTexts.value[comment.ID]) {
        // 检查是否包含代码块或LaTeX公式
        const hasCodeBlock = comment.quoted_text.includes('```')
        const hasLatex = comment.quoted_text.includes('$$') || comment.quoted_text.includes('$')

        // 如果包含代码块或LaTeX公式，使用renderQuotedText渲染
        if (hasCodeBlock || hasLatex) {
          const rendered = await renderQuotedText(comment.quoted_text, 250)
          commentQuotedTexts.value[comment.ID] = rendered
        } else {
          // 否则只将换行符转换为 <br>，转义HTML特殊字符
          const escapedText = comment.quoted_text
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;')
          commentQuotedTexts.value[comment.ID] = escapedText.replace(/\n/g, '<br>')
        }
      }
    })
  await Promise.all(renderPromises)
  // 渲染完成后检测所有评论引用文本是否溢出
  await nextTick()
  allComments
    .filter(comment => comment.quoted_text && comment.quoted_text.trim())
    .forEach(comment => {
      checkCommentQuotedTextOverflow(comment.ID)
    })
}

const handleTextShare = (text) => {
  shareSelectedText.value = text
  // 生成文章 URL
  const id = props.articleId || route.params.id
  const type = props.type || route.params.type || route.query.type || 'blog'
  const baseUrl = window.location.origin
  if (type === 'moment') {
    articleUrl.value = `${baseUrl}/moments/${id}`
  } else if (type === 'blog') {
    articleUrl.value = `${baseUrl}/blog/${id}`
  } else if (type === 'research' || type === 'project') {
    articleUrl.value = `${baseUrl}/blog/${id}?type=${type}`
  } else {
    articleUrl.value = `${baseUrl}/blog/${id}`
  }
  shareCardVisible.value = true
  // 只隐藏菜单，不清除高亮
  textSelectionMenuVisible.value = false
  selectedText.value = ''
  highlightedRange.value = null
}

// 获取文章副标题（用于分享卡片）
const getArticleSubtitle = () => {
  // 如果文章有摘要，使用摘要；否则使用前100个字符
  const abstract = content.value.substring(0, 100)
  return abstract.length < content.value.length ? abstract + '...' : abstract
}

onMounted(async () => {
  // 延迟执行，确保路由参数完全加载
  await nextTick()

  await initializeDetail()

  // 添加滚动监听器（使用防抖优化性能）
  const debouncedHandleScroll = debounceScroll(handleScroll, 16)
  window.addEventListener('scroll', debouncedHandleScroll, { passive: true })

  // 添加文本选择监听器
  document.addEventListener('mouseup', handleTextSelection)
  document.addEventListener('click', (e) => {
    // 如果点击的不是菜单按钮和分享卡片，只隐藏菜单（不清除高亮）
    if (!e.target.closest('.text-selection-menu') && !e.target.closest('.share-card-container')) {
      setTimeout(() => {
        const selection = window.getSelection()
        if (!selection.toString().trim()) {
          // 只隐藏菜单，不清除高亮
          textSelectionMenuVisible.value = false
          selectedText.value = ''
          highlightedRange.value = null
        }
      }, 100)
    }
  })

  // 监听路由变化，切换文章时清除高亮和引用文本
  watch(() => route.params.id, (newId, oldId) => {
    // 只有在真正切换文章时才清除高亮和引用文本
    if (newId !== oldId && oldId) {
      removeHighlight()
      clearQuotedText()
      commentQuotedTexts.value = {} // 清空评论引用文本缓存
    }
  })

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

// keep-alive 激活时重新初始化（处理路由切换但组件未卸载的情况）
onActivated(async () => {
  await nextTick()

  // 保存当前文章ID，用于判断是否切换了文章
  const previousArticleId = route.params.id || props.articleId

  await initializeDetail()

  // 检查是否切换了文章
  const currentArticleId = route.params.id || props.articleId
  if (previousArticleId !== currentArticleId) {
    // 切换了文章，清除高亮和引用文本
    removeHighlight()
    clearQuotedText()
    commentQuotedTexts.value = {} // 清空评论引用文本缓存
  }
  // 如果没有切换文章，高亮会在 restoreHighlights 中恢复

  // 重新初始化目录位置
  setTimeout(() => {
    adjustTocPosition()
    initializeTocScrollPosition()
  }, 300)
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  // 清理滚动监听器
  if (scrollTimeout) {
    clearTimeout(scrollTimeout)
  }
  // 注意：由于使用了防抖函数，这里需要保存函数引用才能正确移除
  // 在实际应用中，如果需要严格清理，应该保存函数引用

  // 清理文本选择监听器
  document.removeEventListener('mouseup', handleTextSelection)

  // 保存高亮状态（在卸载前保存，以便下次访问时恢复）
  saveHighlights()
})

// 格式化评论时间
// 格式化评论内容，保留换行并支持Markdown
const formatCommentContent = async (content) => {
  if (!content) return ''

  // 确保 marked 已加载
  if (!marked) {
    await loadMarkdownLibs()
  }

  if (marked) {
    // 使用 marked 渲染 Markdown（会自动处理换行）
    const html = marked(content, {
      breaks: true, // 将换行符转换为 <br>
      gfm: true,
      headerIds: false,
      mangle: false,
      sanitize: false
    })

    // 使用 DOMPurify 清理 HTML（确保安全）
    const DOMPurify = (await import('dompurify')).default
    return DOMPurify.sanitize(html, {
      ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'code', 'pre', 'a', 'ul', 'ol', 'li'],
      ALLOWED_ATTR: ['href', 'class']
    })
  } else {
    // 如果 marked 未加载，只处理换行符
    return content.replace(/\n/g, '<br>')
  }
}

// 缓存格式化后的评论内容
const formattedCommentContents = ref({})

// 格式化所有评论内容
const formatAllCommentContents = async () => {
  const allComments = getAllCommentsInOrder()
  const formatPromises = allComments
    .filter(comment => comment.content && comment.content.trim())
    .map(async (comment) => {
      if (!formattedCommentContents.value[comment.ID]) {
        const formatted = await formatCommentContent(comment.content)
        formattedCommentContents.value[comment.ID] = formatted
      }
    })
  await Promise.all(formatPromises)
}

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

// 初始化文章详情（可复用的函数）
const initializeDetail = async () => {
  // 重置页面滚动位置到顶部（使用instant避免触发hover）
  window.scrollTo({ top: 0, behavior: 'instant' })

  const id = props.articleId || route.params.id
  // 只有当ID存在且有效时才加载
  if (id && id !== 'undefined') {
    await loadDetail()
    await Promise.all([
      loadComments(),
      loadLikeStatus() // 并行加载评论和点赞状态
    ])
  }
}

// 监听路由变化，重新加载文章和评论
watch(
  () => [route.params.id, route.params.type, props.articleId, props.type],
  async ([newRouteId, newRouteType, newPropsId, newPropsType], [oldRouteId, oldRouteType, oldPropsId, oldPropsType]) => {
    // 如果路由参数或props发生变化，重新加载
    const idChanged = newRouteId !== oldRouteId || newPropsId !== oldPropsId
    const typeChanged = newRouteType !== oldRouteType || newPropsType !== oldPropsType

    if (idChanged || typeChanged) {
      await initializeDetail()
    }
  },
  { immediate: false }
)

// keep-alive 激活时重新初始化（处理路由切换但组件未卸载的情况）
// 监听输入框引用文本的变化，检测是否溢出
watch(renderedQuotedText, async () => {
  await nextTick()
  checkInputQuotedTextOverflow()
}, { flush: 'post' })

// 监听评论引用文本的变化，检测是否溢出
watch(commentQuotedTexts, async () => {
  await nextTick()
  Object.keys(commentQuotedTextRefs.value).forEach(commentId => {
    checkCommentQuotedTextOverflow(Number(commentId))
  })
}, { deep: true, flush: 'post' })

watch(
  () => route.path,
  async (newPath, oldPath) => {
    if (newPath !== oldPath) {
      await initializeDetail()
    }
  },
  { immediate: false }
)

// 为标题添加标识标签
const addHeadingLabels = () => {
  const markdownBody = document.querySelector('.markdown-body')
  if (!markdownBody) return

  // 为每个标题添加标识标签
  const headings = markdownBody.querySelectorAll('h1, h2, h3, h4, h5, h6')
  headings.forEach((heading) => {
    // 如果已经添加过标识，跳过
    if (heading.querySelector('.heading-label')) return

    const level = parseInt(heading.tagName.charAt(1))
    const label = document.createElement('span')
    label.className = 'heading-label'
    label.setAttribute('data-level', level.toString())
    label.textContent = `H${level}`

    // 直接设置内联样式确保样式生效
    label.style.display = 'inline-block'
    label.style.padding = level === 1 ? '3px 10px' : level === 2 ? '2px 9px' : level === 3 ? '2px 8px' : level === 4 ? '2px 7px' : '1px 6px'
    label.style.borderRadius = '4px'
    label.style.marginRight = '10px'
    label.style.fontWeight = '700'
    label.style.fontSize = level === 1 ? '0.65em' : level === 2 || level === 3 ? '0.7em' : level === 4 ? '0.65em' : '0.6em'
    label.style.fontFamily = "'Inter', 'Noto Sans SC', sans-serif"
    label.style.letterSpacing = '0.5px'
    label.style.lineHeight = '1.4'
    label.style.boxShadow = '0 2px 4px rgba(139, 92, 246, 0.3)'
    label.style.flexShrink = '0'
    label.style.border = 'none'
    label.style.textDecoration = 'none'

    // 根据级别设置不同的背景色和文字颜色
    if (level === 1) {
      label.style.background = 'linear-gradient(135deg, #7c3aed 0%, #8b5cf6 100%)'
      label.style.color = 'white'
    } else if (level === 2) {
      label.style.background = 'linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%)'
      label.style.color = 'white'
    } else if (level === 3) {
      label.style.background = 'linear-gradient(135deg, #a855f7 0%, #c084fc 100%)'
      label.style.color = 'white'
    } else if (level === 4) {
      label.style.background = 'linear-gradient(135deg, #c084fc 0%, #d8b4fe 100%)'
      label.style.color = 'white'
    } else if (level === 5) {
      label.style.background = 'linear-gradient(135deg, #d8b4fe 0%, #e9d5ff 100%)'
      label.style.color = '#7c3aed'
    } else if (level === 6) {
      label.style.background = 'linear-gradient(135deg, #e9d5ff 0%, #f3e8ff 100%)'
      label.style.color = '#8b5cf6'
    }

    // 将标识插入到标题的最前面
    if (heading.firstChild) {
      heading.insertBefore(label, heading.firstChild)
    } else {
      heading.appendChild(label)
    }
  })
}

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
  flex-wrap: nowrap; /* 强制保持一行 */
}

.left-buttons, .right-buttons {
  display: flex;
  gap: 15px;
}

.like-btn, .subscribe-btn, .comment-btn, .share-btn, .edit-btn {
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
  white-space: nowrap; /* 防止文字换行 */
  flex-shrink: 0; /* 默认不收缩，但在小屏幕时可以收缩 */
}

.like-btn:hover, .subscribe-btn:hover, .comment-btn:hover, .share-btn:hover, .edit-btn:hover {
  background: #f8f9fa;
  border-color: #d0d0d0;
  color: #333;
}

.like-btn svg {
  color: #ff6b6b;
  transition: all 0.3s ease;
}

.like-btn.liked svg {
  color: #ff4757;
  animation: likeAnimation 0.5s ease;
}

@keyframes likeAnimation {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
  }
}

.like-btn.liked {
  border-color: #ff4757;
  background: rgba(255, 71, 87, 0.1);
}

.like-btn.loading {
  opacity: 0.6;
  cursor: not-allowed;
  pointer-events: none;
}

.like-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.edit-btn svg {
  color: #a855f7;
}

.edit-btn:hover svg {
  color: #7c3aed;
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

/* 确保 markdown-body 中的 katex-block 不受父元素影响 */
.markdown-body .katex-block,
.markdown-body p .katex-block,
.markdown-body div .katex-block {
  text-align: center !important;
  display: block !important;
  width: 100% !important;
  margin: 24px 0 !important;
  max-width: 100% !important;
}

/* 确保包裹 katex-block 的父元素不影响居中 */
.markdown-body p:has(.katex-block),
.markdown-body div:has(.katex-block) {
  text-align: center !important;
}

/* 文章内容中的图片样式 - 宽度不超过容器80%，高度不超过750px，居中显示 */
.article-content .markdown-body img,
.content-container .markdown-body img,
.detail-view .markdown-body img,
.markdown-body img {
  max-width: 80% !important;
  max-height: 750px !important;
  width: auto !important;
  height: auto !important;
  display: block !important;
  margin: 20px auto !important;
  border-radius: 8px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
  object-fit: contain !important;
  box-sizing: border-box !important;
}

/* 先定义 figcaption 和 LaTeX 公式的样式，确保优先级 */
.markdown-body figure figcaption,
.markdown-body figcaption {
  text-align: center !important;
  display: block !important;
}

.markdown-body *:not(figcaption):not(.katex-block):not(.katex) {
  text-align: left !important;
}

/* 标题标识标签样式 - 使用 !important 确保优先级 */
.markdown-body h1 .heading-label,
.markdown-body h2 .heading-label,
.markdown-body h3 .heading-label,
.markdown-body h4 .heading-label,
.markdown-body h5 .heading-label,
.markdown-body h6 .heading-label {
  display: inline-block !important;
  background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%) !important;
  color: white !important;
  font-size: 0.7em !important;
  font-weight: 700 !important;
  padding: 2px 8px !important;
  border-radius: 4px !important;
  margin-right: 10px !important;
  vertical-align: middle !important;
  line-height: 1.4 !important;
  font-family: 'Inter', 'Noto Sans SC', sans-serif !important;
  letter-spacing: 0.5px !important;
  box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3) !important;
  flex-shrink: 0 !important;
  border: none !important;
  text-decoration: none !important;
}

/* 不同级别标题的标识样式 */
.markdown-body h1 .heading-label {
  font-size: 0.65em !important;
  padding: 3px 10px !important;
  background: linear-gradient(135deg, #7c3aed 0%, #8b5cf6 100%) !important;
}

.markdown-body h2 .heading-label {
  font-size: 0.7em !important;
  padding: 2px 9px !important;
  background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%) !important;
}

.markdown-body h3 .heading-label {
  font-size: 0.7em !important;
  padding: 2px 8px !important;
  background: linear-gradient(135deg, #a855f7 0%, #c084fc 100%) !important;
}

.markdown-body h4 .heading-label {
  font-size: 0.65em !important;
  padding: 2px 7px !important;
  background: linear-gradient(135deg, #c084fc 0%, #d8b4fe 100%) !important;
}

.markdown-body h5 .heading-label {
  font-size: 0.6em !important;
  padding: 1px 6px !important;
  background: linear-gradient(135deg, #d8b4fe 0%, #e9d5ff 100%) !important;
  color: #7c3aed !important;
}

.markdown-body h6 .heading-label {
  font-size: 0.6em !important;
  padding: 1px 6px !important;
  background: linear-gradient(135deg, #e9d5ff 0%, #f3e8ff 100%) !important;
  color: #8b5cf6 !important;
}

/* 确保标题和标识在同一行显示 */
.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  text-align: left !important;
  display: flex !important;
  align-items: center !important;
  flex-wrap: wrap !important;
  gap: 0 !important;
}

/* 标题内容部分（排除标识） */
.markdown-body h1 > *:not(.heading-label),
.markdown-body h2 > *:not(.heading-label),
.markdown-body h3 > *:not(.heading-label),
.markdown-body h4 > *:not(.heading-label),
.markdown-body h5 > *:not(.heading-label),
.markdown-body h6 > *:not(.heading-label) {
  flex: 1;
  min-width: 0;
}

.markdown-body p,
.markdown-body div:not(.katex-block),
.markdown-body span:not(.katex) {
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

/* LaTeX 公式样式 - 必须在通用样式之后，确保优先级 */
.markdown-body .katex-block,
.markdown-body p .katex-block,
.markdown-body div .katex-block {
  margin: 24px 0 !important;
  text-align: center !important;
  overflow-x: auto;
  overflow-y: hidden;
  display: block !important;
  width: 100% !important;
  max-width: 100% !important;
}

.markdown-body .katex-block .katex,
.markdown-body p .katex-block .katex,
.markdown-body div .katex-block .katex {
  font-size: 1.2em;
  display: inline-block !important;
  text-align: center !important;
  margin: 0 auto !important;
}

/* 确保 katex 元素本身也居中 */
.markdown-body .katex {
  text-align: center !important;
}

.markdown-body .katex {
  font-size: 1.05em;
  line-height: 1.8;
}

/* 优化上标和下标间距 - 增加上标与基线的距离 */
.markdown-body .katex .msup {
  margin-left: 0.15em;
}

.markdown-body .katex .msub {
  margin-left: 0.15em;
}

/* 优化上标内部间距 */
.markdown-body .katex .msup > .vlist-t {
  margin-top: 0.1em;
}

/* 优化运算符间距 */
.markdown-body .katex .mop {
  margin-left: 0.16667em;
  margin-right: 0.16667em;
}

.markdown-body .katex .mord + .mop {
  margin-left: 0.16667em;
}

.markdown-body .katex .mop + .mord {
  margin-left: 0.16667em;
}

.markdown-body .katex-error {
  color: #d1242f;
  background-color: #fff5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
}

/* 调试样式已移除 */

/* 评论区 */
.comments-section {
  background: transparent;
  border-radius: 0;
  padding: 30px 0;
  box-shadow: none;
  border: none;
  margin-bottom: 30px;
}

.section-header {
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(139, 92, 246, 0.1);
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

.comment-icon {
  width: 24px;
  height: 24px;
  color: #8b5cf6;
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
  padding: 0;
  margin-bottom: 2rem;
  transition: all 0.3s ease;
  position: relative;
  align-items: flex-start;
}

.comment-item:hover {
  transform: translateY(-1px);
}

.comment-avatar {
  flex-shrink: 0;
}

.avatar-circle {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #00BFFF;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  font-weight: 600;
  text-transform: uppercase;
  box-shadow: 0 2px 8px rgba(0, 191, 255, 0.3);
  transition: transform 0.2s ease;
}

.avatar-circle:hover {
  transform: scale(1.05);
}

.comment-body {
  flex: 1;
  min-width: 0;
  max-width: calc(100% - 60px);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.comment-bubble {
  background: white;
  border-radius: 12px;
  padding: 0.85rem 1.05rem;
  border-top: 1.5px solid rgba(139, 92, 246, 0.4);
  border-left: 1.5px solid rgba(139, 92, 246, 0.4);
  border-bottom: 1px solid #f0f0f0;
  border-right: 1px solid #f0f0f0;
  box-shadow:
    0 4px 12px rgba(139, 92, 246, 0.2),
    0 2px 6px rgba(139, 92, 246, 0.15),
    0 1px 2px rgba(0, 0, 0, 0.05);
  position: relative;
  word-wrap: break-word;
  display: inline-block;
  max-width: 100%;
  margin-top: 0.3rem;
  transition: box-shadow 0.2s ease, border-color 0.2s ease;
}

/* 箭头外层边框（包含顶部紫色边框） */
.comment-bubble::after {
  content: '';
  position: absolute;
  left: -9px;
  top: 11px;
  width: 0;
  height: 0;
  border-top: 9px solid transparent;
  border-bottom: 9px solid transparent;
  border-right: 9px solid rgba(139, 92, 246, 0.4);
  z-index: 0;
}

/* 箭头内层白色填充 */
.comment-bubble::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 12px;
  width: 0;
  height: 0;
  border-top: 8px solid transparent;
  border-bottom: 8px solid transparent;
  border-right: 8px solid white;
  z-index: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  position: relative;
}

.hover-reply-btn {
  opacity: 0;
  padding: 0.2rem 0.5rem;
  background: #8b5cf6;
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-left: auto;
}

.comment-item:hover .hover-reply-btn {
  opacity: 1;
}

.comment-body.selected {
  background: rgba(139, 92, 246, 0.1);
  border-radius: 8px;
  padding: 0.5rem;
  margin: -0.5rem;
}

.reply-tag {
  background: #8b5cf6;
  color: white;
  padding: 0.15rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
  margin: 0 0.3rem;
  display: inline-block;
  white-space: nowrap;
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
  color: #2b2b2b;
  line-height: 1.6;
  word-wrap: break-word;
  font-size: 0.95rem;
  margin: 0;
  text-align: left;
}

/* 内联代码样式 */
.comment-content :deep(code) {
  background: #f1f5f9;
  color: #475569;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 0.85em;
  border: 1px solid #e2e8f0;
  font-weight: 500;
}

/* 代码块样式 */
.comment-content :deep(pre) {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem;
  margin: 0.5rem 0;
  overflow-x: auto;
}

.comment-content :deep(pre code) {
  background: none;
  border: none;
  padding: 0;
  color: #1e293b;
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

/* 引用文本容器 */
.quoted-text-container {
  position: relative;
  background: rgba(168, 85, 247, 0.15);
  border-left: 3px solid rgba(168, 85, 247, 0.5);
  border-radius: 8px;
  padding: 12px 16px;
  margin: 8px;
  margin-bottom: 8px;
  font-size: 0.9rem;
  line-height: 1.6;
}

.quoted-text-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #7c3aed;
  margin-bottom: 8px;
  text-align: left;
}

.quoted-text-content {
  color: #333;
  max-height: 7.2em; /* 约6行的高度 (1.2 * 6) */
  overflow: hidden;
  background: transparent !important;
  word-break: break-word;
  line-height: 1.2;
  position: relative;
  padding-right: 48px; /* 提前两个字截断，为省略号留出空间（约2个中文字符宽度） */
}

/* 确保所有子元素也遵守截断规则 */
.quoted-text-content :deep(*) {
  max-width: 100%;
  overflow: hidden;
}

/* 添加省略号（当内容被截断时） */
.quoted-text-content.is-truncated::after {
  content: '...';
  position: absolute;
  right: 4px;
  bottom: 0;
  background: transparent; /* 透明背景，因为容器已有背景 */
  padding-left: 8px;
  padding-right: 4px;
  pointer-events: none;
  z-index: 2;
}

.quoted-text-content.markdown-body {
  background: transparent !important;
}

.quoted-text-content.markdown-body :deep(*),
.quoted-text-content.markdown-body :deep(p),
.quoted-text-content.markdown-body :deep(div),
.quoted-text-content.markdown-body :deep(span),
.quoted-text-content.markdown-body :deep(code),
.quoted-text-content.markdown-body :deep(pre) {
  background: transparent !important;
}

.quoted-text-content :deep(p),
.quoted-text-content :deep(div) {
  margin: 0.2em 0;
  display: block;
  line-height: 1.2;
}

.quoted-text-content :deep(br) {
  display: block;
  line-height: 1.2;
  height: 1.2em;
}

.quoted-text-content :deep(p:first-child),
.quoted-text-content :deep(div:first-child) {
  margin-top: 0;
}

.quoted-text-content :deep(p:last-child),
.quoted-text-content :deep(div:last-child) {
  margin-bottom: 0;
}

.cancel-quote-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  border: none;
  background: rgba(168, 85, 247, 0.2);
  color: #7c3aed;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  font-size: 0.75rem;
}

.cancel-quote-btn:hover {
  background: rgba(168, 85, 247, 0.3);
  transform: scale(1.1);
}

/* 评论中的引用文本 */
.comment-quoted-text {
  background: rgba(168, 85, 247, 0.15);
  border-left: 3px solid rgba(168, 85, 247, 0.5);
  border-radius: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
  font-size: 0.85rem;
  line-height: 1.5;
  margin-top: 0;
}

.comment-quoted-text .quoted-text-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #7c3aed;
  margin-bottom: 6px;
  text-align: left;
}

.comment-quoted-text .quoted-text-content {
  color: #333;
  max-height: 6em; /* 约5行的高度 (1.2 * 5) */
  overflow: hidden;
  background: transparent !important;
  word-break: break-word;
  line-height: 1.2;
  position: relative;
  padding-right: 48px; /* 提前两个字截断，为省略号留出空间（约2个中文字符宽度） */
}

/* 确保所有子元素也遵守截断规则 */
.comment-quoted-text .quoted-text-content :deep(*) {
  max-width: 100%;
  overflow: hidden;
}

/* 添加省略号（当内容被截断时） */
.comment-quoted-text .quoted-text-content.is-truncated::after {
  content: '...';
  position: absolute;
  right: 4px;
  bottom: 0;
  background: transparent; /* 透明背景，因为容器已有背景 */
  padding-left: 8px;
  padding-right: 4px;
  pointer-events: none;
  z-index: 2;
}

.comment-quoted-text .quoted-text-content.markdown-body {
  background: transparent !important;
}

.comment-quoted-text .quoted-text-content.markdown-body :deep(*),
.comment-quoted-text .quoted-text-content.markdown-body :deep(p),
.comment-quoted-text .quoted-text-content.markdown-body :deep(div),
.comment-quoted-text .quoted-text-content.markdown-body :deep(span),
.comment-quoted-text .quoted-text-content.markdown-body :deep(code),
.comment-quoted-text .quoted-text-content.markdown-body :deep(pre) {
  background: transparent !important;
}

.comment-quoted-text .quoted-text-content :deep(p),
.comment-quoted-text .quoted-text-content :deep(div) {
  margin: 0.2em 0;
  font-size: 0.9em;
  display: block;
  line-height: 1.2;
}

.comment-quoted-text .quoted-text-content :deep(br) {
  display: block;
  line-height: 1.2;
  height: 1.2em;
}

.comment-quoted-text .quoted-text-content :deep(p:first-child),
.comment-quoted-text .quoted-text-content :deep(div:first-child) {
  margin-top: 0;
}

.comment-quoted-text .quoted-text-content :deep(p:last-child),
.comment-quoted-text .quoted-text-content :deep(div:last-child) {
  margin-bottom: 0;
}

/* 评论编辑器 */
.comment-editor {
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid rgba(139, 92, 246, 0.1);
  background: transparent;
  border-radius: 0;
  padding: 30px 0;
  box-shadow: none;
}

.comment-editor h3 {
  font-size: 1.3rem;
  color: #333;
  margin-bottom: 20px;
}

.editor-wrapper {
  background: transparent;
  border-radius: 12px;
  padding: 0;
  border: none;
  transition: all 0.3s ease;
}

.editor-wrapper:focus-within {
  border-bottom-color: rgba(102, 126, 234, 0.3);
  background: transparent; /* 仅输入框变色，容器不变色 */
  box-shadow: none;
  backdrop-filter: none;
}

.comment-input-wrapper {
  position: relative;
  width: 100%;
}

.comment-input-box {
  position: relative;
  width: 100%;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  background: white;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
  min-height: 120px;
}

.comment-input-box:focus-within {
  border-color: rgba(168, 85, 247, 0.5);
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.1);
}

.comment-input {
  width: 100%;
  min-height: 120px;
  padding: 15px 50px 15px 15px;
  border: none;
  border-radius: 0;
  font-size: 0.95rem;
  line-height: 1.6;
  resize: vertical;
  background: transparent;
  color: #333;
  font-family: inherit;
  transition: none;
}

.comment-input-box.has-quoted-text .comment-input {
  border-radius: 0 0 12px 12px;
}

.emoji-btn {
  position: absolute;
  right: 12px;
  bottom: 12px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s;
  line-height: 1;
  font-family: "Apple Color Emoji", "Segoe UI Emoji", "Noto Color Emoji", "Segoe UI Symbol", "Android Emoji", "EmojiSymbols", "EmojiOne Mozilla", "Twemoji Mozilla", "Segoe UI", sans-serif;
}

.emoji-btn:hover {
  background: #f3f4f6;
  transform: scale(1.1);
}

.emoji-btn:active {
  transform: scale(0.95);
}

.comment-input:focus {
  outline: none;
}

.comment-input::placeholder {
  color: #9ca3af;
}

.editor-actions {
  display: flex;
  align-items: center;
  margin-top: 15px;
}

.reply-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #8b5cf6;
  border-radius: 8px;
  color: white;
}

.reply-label {
  font-size: 0.9rem;
  font-weight: 500;
}

.cancel-reply-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 4px;
  color: white;
  cursor: pointer;
  padding: 0.2rem 0.4rem;
  font-size: 0.8rem;
  transition: background 0.3s ease;
}

.cancel-reply-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: auto;
}

.char-count {
  color: #999;
  font-size: 0.9rem;
}

.preview-btn {
  padding: 12px 24px;
  background: #8b5cf6;
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3);
}

.preview-btn:hover {
  background: #7c3aed;
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
  tab-size: 4 !important;
  -moz-tab-size: 4 !important;
}

/* 确保代码块中的 Tab 显示为 4 个空格宽度 */
.markdown-body pre code {
  tab-size: 4 !important;
  -moz-tab-size: 4 !important;
}

/* 文章目录样式 - 初始隐藏，等待JS定位 */
.floating-toc {
  position: fixed;
  top: 80px;
  left: 20px; /* 向屏幕左边靠近 */
  width: 300px;
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
  overflow: hidden;
  width: 100%;
  text-align: left; /* 强制左对齐 */
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
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
  text-align: left; /* 强制左对齐 */
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
    width: 210px;
    /* 定位完全由JS控制 */
  }
}

/* 大屏幕优化 - 完全由JS控制 */

/* 超大屏幕优化 */
@media (min-width: 1600px) {
  .floating-toc {
    width: 240px;
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

/* 文本选择颜色（浏览器默认选中效果） */
.article-content ::selection,
.markdown-body ::selection {
  background-color: rgba(124, 58, 237, 0.8);
  color: #ffffff;
}

.article-content ::-moz-selection,
.markdown-body ::-moz-selection {
  background-color: rgba(124, 58, 237, 0.8);
  color: #ffffff;
}

/* 文本高亮样式（点击高亮按钮后的效果） */
/* 高亮文本样式 - 使用 span 标签，display: inline 避免破坏布局 */
/* 覆盖浏览器和第三方CSS的默认黄色背景，强制使用紫色 */
.article-content .text-highlight,
.markdown-body .text-highlight,
.text-highlight {
  background-color: rgba(168, 85, 247, 0.4) !important;
  background: rgba(168, 85, 247, 0.4) !important;
  color: inherit !important;
  padding: 0 !important;
  margin: 0 !important;
  border-radius: 0 !important;
  box-shadow: none !important;
  display: inline !important;
  position: relative !important;
  z-index: 1 !important;
  line-height: inherit !important;
  vertical-align: baseline !important;
  white-space: normal !important;
  word-break: normal !important;
}

/* 确保高亮在列表中的样式正确 */
.markdown-body li .text-highlight,
.article-content li .text-highlight {
  display: inline !important;
  line-height: inherit !important;
  vertical-align: baseline !important;
}

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
  background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3);
}

/* 回复评论样式 - 缩进一个头像宽度 */
.reply-comment {
  margin-left: 60px; /* 一个头像宽度 + 间距 */
}

.reply-comment .comment-avatar {
  margin-top: 0.1rem; /* 稍微调整头像位置 */
}

.preview-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #8b5cf6;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.preview-item {
  opacity: 0.8;
  transform: scale(0.98);
  transition: all 0.3s ease;
}

.preview-item:hover {
  opacity: 1;
  transform: scale(1);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(139, 92, 246, 0.4);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  box-shadow: none;
}

/* 响应式 - 平板和中等屏幕优化 */
@media (max-width: 800px) {
  .left-buttons, .right-buttons {
    /* 减少按钮间距，节省空间 */
    gap: 10px;
  }

  .like-btn, .subscribe-btn, .comment-btn, .share-btn, .edit-btn {
    /* 稍微减小按钮 padding，节省空间 */
    padding: 7px 14px;
    font-size: 0.85rem;
  }
}

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
    /* 保持一行布局，不改为垂直布局 */
    margin-bottom: 15px;
    /* space-between 会自动处理间距，允许中间收缩 */
  }

  .left-buttons, .right-buttons {
    gap: 8px; /* 进一步减小间距 */
  }

  .like-btn, .subscribe-btn, .comment-btn, .share-btn, .edit-btn {
    padding: 6px 12px;
    font-size: 0.8rem;
    min-width: 0; /* 允许按钮内容收缩 */
    flex-shrink: 1; /* 允许按钮收缩 */
  }

  .article-stats {
    font-size: 0.8rem;
  }

  .article-image {
    max-height: 300px;
    /* 确保图片在移动端也有侧边距 */
    margin: 0 20px;
    max-width: calc(100% - 40px); /* 减去左右各20px的边距 */
  }

  /* 确保内容中的图片在移动端也保持侧边距 */
  .article-content .markdown-body img,
  .content-container .markdown-body img,
  .detail-view .markdown-body img,
  .markdown-body img {
    max-width: calc(100% - 40px) !important; /* 减去左右各20px的边距，确保至少有20px侧边距 */
    margin: 20px auto !important;
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
