<template>
  <div class="media-card-wrapper">
    <div class="media-card" :class="{ 'media-row': props.variant === 'row' }" @click="openMediaModal">
      <div v-if="props.variant !== 'row'" class="poster" :style="{ backgroundImage: `url(${props.media.Poster})` }"></div>
      <div class="content">
        <template v-if="props.variant === 'row'">
          <div class="row-header">
            <h2 class="row-title">
              <font-awesome-icon :icon="mediaTypeIcon" class="type-fa-icon" />
              {{ props.media.Name }}
            </h2>
            <div v-if="props.media.Rating" class="row-rating">
              <span class="stars">
                <span v-for="(s, idx) in starArray" :key="'s-'+idx" class="star-wrap">
                  <span v-if="s.cls === 'half-star'" class="star-half-wrap">
                    <font-awesome-icon :icon="['far','star']" class="star-icon empty-star" />
                    <font-awesome-icon :icon="'star'" class="star-icon full-star star-half-overlay" />
                  </span>
                  <font-awesome-icon v-else-if="s.cls === 'full-star'" :icon="s.icon" class="star-icon full-star" />
                  <font-awesome-icon v-else :icon="s.icon" class="star-icon empty-star" />
                </span>
              </span>
              <span class="rating-text">{{ props.media.Rating }}/10</span>
            </div>
          </div>
          <div class="row-meta">
            <span class="date">Modify {{ formattedTime }}</span>
          </div>
        </template>
        <template v-else>
          <h2>
            <font-awesome-icon :icon="mediaTypeIcon" class="type-fa-icon" />
            {{ props.media.Name }}
          </h2>
          <p class="date">Modify {{ formattedTime }}</p>
          <div v-if="props.media.Rating" class="rating-stars">
            <span class="stars">
              <span v-for="(s, idx) in starArray" :key="'m-'+idx" class="star-wrap">
                <span v-if="s.cls === 'half-star'" class="star-half-wrap">
                  <font-awesome-icon :icon="['far','star']" class="star-icon empty-star" />
                  <font-awesome-icon :icon="'star'" class="star-icon full-star star-half-overlay" />
                </span>
                <font-awesome-icon v-else-if="s.cls === 'full-star'" :icon="s.icon" class="star-icon full-star" />
                <font-awesome-icon v-else :icon="s.icon" class="star-icon empty-star" />
              </span>
            </span>
          </div>
        </template>
      </div>
    </div>

    <!-- 媒体详情模态框 - 与搜索模态风格一致 -->
    <Teleport to="body">
      <transition name="media-modal-fade">
        <div v-if="showMediaModal" class="media-modal-overlay" @click="closeMediaModal">
          <div class="media-modal-container" @click.stop>
            <div class="media-modal-header">
              <div v-if="props.media.Poster" class="modal-thumb">
                <img :src="props.media.Poster" :alt="props.media.Name" loading="lazy" decoding="async" @error="onImgError($event)" />
              </div>
              <div class="media-modal-headinfo">
                <h3 class="media-modal-title">{{ props.media.Name }}</h3>
                <div v-if="props.media.Rating" class="media-modal-rating">
                  <span class="stars modal-stars">
                    <span v-for="(s, idx) in starArray" :key="'d-'+idx" class="star-wrap">
                      <span v-if="s.cls === 'half-star'" class="star-half-wrap">
                        <font-awesome-icon :icon="['far','star']" class="star-icon empty-star" />
                        <font-awesome-icon :icon="'star'" class="star-icon full-star star-half-overlay" />
                      </span>
                      <font-awesome-icon v-else-if="s.cls === 'full-star'" :icon="s.icon" class="star-icon full-star" />
                      <font-awesome-icon v-else :icon="s.icon" class="star-icon empty-star" />
                    </span>
                  </span>
                  <span class="rating-text">{{ props.media.Rating }}/10</span>
                </div>
              </div>
              <button v-if="userLevel >= 3" class="media-modal-edit" @click="goToEdit">编辑</button>
            </div>
            <div class="media-modal-body">
              <div class="media-modal-content">
                <!-- eslint-disable-next-line vue/no-v-html -->
                <div class="media-modal-review markdown-body" v-html="renderedReview"></div>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { showCustomMessage } from '@/utils/waifuMessage'
import { protectLatex, restoreAndRenderLatex } from '@/utils/latex'

const props = defineProps({
  media: Object,
  type: String,
  variant: { type: String, default: 'card' } // 'card' | 'row'
})
const mediaTypeIcon = computed(() => {
  const t = (props.type || '').toLowerCase()
  if (t === 'books') return 'book'
  if (t === 'novels') return 'bookmark'
  if (t === 'movies') return 'film'
  return 'box'
})

// 无对外事件

const store = useStore()
const router = useRouter()
const user = store.state.user
const userLevel = user.level

// 取消卡片悬停展示内容，仅保留普通悬停效果
const showMediaModal = ref(false)

// 格式化时间
const formattedTime = computed(() => {
  const date = new Date(props.media.CreatedAt)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
})

// 评分转为 5 颗星（每星2分）：满/半/空
const starArray = computed(() => {
  const rating = Number(props.media.Rating || 0)
  const stars = Math.max(0, Math.min(5, rating / 2))
  const full = Math.floor(stars)
  const hasHalf = stars - full >= 0.5 ? 1 : 0
  const empty = 5 - full - hasHalf
  const arr = []
  for (let i = 0; i < full; i++) arr.push({ icon: 'star', cls: 'full-star' })
  if (hasHalf) arr.push({ icon: 'star-half-stroke', cls: 'half-star' })
  for (let i = 0; i < empty; i++) arr.push({ icon: ['far', 'star'], cls: 'empty-star' })
  return arr
})

// 渲染 Markdown（使用 ref 因为需要异步处理 LaTeX）
const renderedReview = ref('')

// 异步渲染 Markdown 内容
const renderReviewContent = async () => {
  if (!props.media.Review) {
    renderedReview.value = ''
    return
  }

  try {
    // 预处理：保护 LaTeX 公式（在 marked 渲染前处理）
    const { protected: protectedContent, placeholders: latexPlaceholders } = protectLatex(props.media.Review)
    let html = marked(protectedContent, {
      breaks: true, // 将换行符转换为 <br>
      gfm: true, // 启用 GitHub Flavored Markdown
      headerIds: false,
      mangle: false,
      sanitize: false
    })

    // 为图片添加内联样式，确保在渲染时就有宽度和高度限制，避免闪烁
    html = html.replace(
      /<img([^>]*)(style="[^"]*")?([^>]*)>/gi,
      (match, before, existingStyle, after) => {
        // 如果已经有style属性，则合并样式
        if (existingStyle) {
          const newStyle = existingStyle.replace(/"/g, '') + '; max-width: 80% !important; max-height: 400px !important; width: auto !important; height: auto !important; display: block !important; margin: 20px auto !important; box-sizing: border-box !important;'
          return `<img${before} style="${newStyle}"${after}>`
        } else {
          // 如果没有style属性，则添加新的style属性
          return `<img${before} style="max-width: 80% !important; max-height: 400px !important; width: auto !important; height: auto !important; display: block !important; margin: 20px auto !important; box-sizing: border-box !important;"${after}>`
        }
      }
    )

    // 恢复并渲染 LaTeX 公式（异步）
    html = await restoreAndRenderLatex(html, latexPlaceholders)

    // 配置DOMPurify允许style属性，确保内联样式不被过滤，支持 LaTeX 公式
    renderedReview.value = DOMPurify.sanitize(html, {
      ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'span', 'div', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'img', 'math', 'annotation', 'semantics', 'mrow', 'mi', 'mo', 'mn', 'mfrac', 'msup', 'msub', 'munderover', 'munder', 'mover', 'mtable', 'mtr', 'mtd', 'mtext'],
      ALLOWED_ATTR: ['class', 'style', 'src', 'alt', 'title', 'width', 'height', 'data-*', 'aria-*', 'role', 'dir', 'colspan', 'rowspan'],
      ALLOW_DATA_ATTR: true
    })
  } catch (error) {
    console.error('渲染评论内容失败:', error)
    renderedReview.value = props.media.Review || ''
  }
}

// 监听 media.Review 变化，重新渲染
watch(() => props.media.Review, () => {
  renderReviewContent()
}, { immediate: true })

const openMediaModal = () => {
  showMediaModal.value = true
  // 确保内容已渲染
  renderReviewContent()
  // 在弹框打开后，延迟为图片添加样式（确保 DOM 完全渲染）
  nextTick(() => {
    setTimeout(() => {
      const modalReview = document.querySelector('.media-modal-review')
      if (modalReview) {
        const images = modalReview.querySelectorAll('img')
        images.forEach(img => {
          // 确保图片样式正确应用（内联样式优先级最高）
          img.style.setProperty('max-width', '80%', 'important')
          img.style.setProperty('max-height', '400px', 'important')
          img.style.setProperty('width', 'auto', 'important')
          img.style.setProperty('height', 'auto', 'important')
          img.style.setProperty('display', 'block', 'important')
          img.style.setProperty('margin', '20px auto', 'important')
          img.style.setProperty('box-sizing', 'border-box', 'important')
        })
      }
    }, 300) // 增加延迟时间，确保 v-html 渲染完成
  })
  // 以 1/4 概率在弹框打开时根据类型说一句话
  try {
    if (Math.random() < 0.25) {
      const t = (props.type || '').toLowerCase()
      let texts = []
      if (t === 'books') {
        texts = [
          '您有多久没有完整读过一本书？那里面有与我们共鸣的灵魂',
          '我记不清上一次打开非专业的纸质书是什么时候了，但它们一定组成了我的一部分',
          '也许您有好书要推荐给我，朋友？',
          '我会永远怀念小时候和亲人朋友在书店里面打发的时光，哪怕那时候的书店朴实的只有书'
        ]
      } else if (t === 'novels') {
        texts = [
          '仔细回想，小说已经贯穿了我到现在为止超过一半的人生',
          '虚构的故事里藏着最真实的情感，欢迎来到想象的世界',
          '我们在小说中寻找想象里的自己，但不要忘了自己不是找到的，而是创造的'
        ]
      } else if (t === 'movies') {
        texts = [
          '无关乎题材与风格，我总是痴迷于所有的好电影',
          '我喜欢有年代感的电影，它里面有时代的影子',
          '好的电影不厌百回看，也许现在是时候翻出你最爱的电影了'
        ]
      }
      if (texts.length) {
        const text = texts[Math.floor(Math.random() * texts.length)]
        showCustomMessage(text, 4000)
      }
    }
  } catch (e) { /* 忽略提示失败 */ }

  // 弹框打开后，为正文中的图片添加错误处理
  // 使用 setTimeout 确保 DOM 完全渲染
  setTimeout(() => {
    const modalReview = document.querySelector('.media-modal-review')
    if (modalReview) {
      const images = modalReview.querySelectorAll('img')
      images.forEach(img => {
        // 避免重复添加
        if (!img.dataset.errorHandlerAdded) {
          img.dataset.errorHandlerAdded = 'true'
          img.addEventListener('error', () => {
            img.style.display = 'none'
          })
        }
      })
    }
  }, 100)
}

// 监听弹框打开，为正文中的图片添加错误处理和样式
watch(showMediaModal, (isOpen) => {
  if (isOpen) {
    // 确保内容已渲染
    renderReviewContent()
    // 使用 nextTick 和延迟执行，确保 v-html 渲染完成后再处理图片
    nextTick(() => {
      setTimeout(() => {
        const modalReview = document.querySelector('.media-modal-review')
        if (modalReview) {
          const images = modalReview.querySelectorAll('img')
          images.forEach(img => {
            // 避免重复添加
            if (!img.dataset.errorHandlerAdded) {
              img.dataset.errorHandlerAdded = 'true'
              img.addEventListener('error', () => {
                img.style.display = 'none'
              })
              // 确保图片样式正确应用（内联样式优先级最高）
              img.style.setProperty('max-width', '80%', 'important')
              img.style.setProperty('max-height', '400px', 'important')
              img.style.setProperty('width', 'auto', 'important')
              img.style.setProperty('height', 'auto', 'important')
              img.style.setProperty('display', 'block', 'important')
              img.style.setProperty('margin', '20px auto', 'important')
              img.style.setProperty('box-sizing', 'border-box', 'important')
            } else {
              // 即使已经添加过错误处理，也要确保样式应用
              img.style.setProperty('max-width', '80%', 'important')
              img.style.setProperty('max-height', '400px', 'important')
              img.style.setProperty('width', 'auto', 'important')
              img.style.setProperty('height', 'auto', 'important')
              img.style.setProperty('display', 'block', 'important')
              img.style.setProperty('margin', '20px auto', 'important')
              img.style.setProperty('box-sizing', 'border-box', 'important')
            }
          })
        }
      }, 300) // 增加延迟时间，确保 v-html 和图片加载完成
    })
  }
})

// 图片错误处理：隐藏图片而不是显示默认图片
const onImgError = (e) => {
  const img = e?.target
  if (img) {
    // 隐藏图片元素
    img.style.display = 'none'
  }
}

const closeMediaModal = () => {
  showMediaModal.value = false
}

// ESC 关闭
const handleKeydown = (e) => {
  if (e.key === 'Escape') closeMediaModal()
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})

// 跳转到编辑页面
const goToEdit = () => {
  // 先关闭模态框
  closeMediaModal()
  router.push({
    path: `/edit/${props.media.ID}`,
    query: {
      contentType: 'media',
      mediaType: props.type
    }
  })
}

// 删除功能已移除（按需可恢复）

// 暴露方法供父组件调用
defineExpose({
  openMediaModal
})
</script>

<style>
/* 导入 Markdown 样式（全局） */
@import '@/assets/styles/github-markdown.css';
</style>

<style scoped>
.media-card {
  position: relative;
  cursor: pointer;
  border-radius: 0;
  overflow: visible;
  background: transparent; /* 融入背景 */
  box-shadow: none;
  height: 450px;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.media-card.media-row {
  display: flex;
  align-items: center;
  height: auto;
  min-height: 84px;
  padding: 14px 18px;
  box-shadow: none;
  border-radius: 0;
  background: transparent;
}

.media-card:hover {
  transform: none;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.12);
}

/* 流光边框效果 */
.media-card::before,
.media-card::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.25s ease;
  border-radius: inherit;
}

.media-card::before {
  /* 顶边（自左向右） + 右边（自上而下） */
  background:
    linear-gradient(90deg, rgba(168,85,247,0.9), rgba(168,85,247,0.9)) top left / 0% 2px no-repeat,
    linear-gradient(180deg, rgba(124,58,237,0.9), rgba(124,58,237,0.9)) top right / 2px 0% no-repeat;
}

.media-card::after {
  /* 底边（自右向左） + 左边（自下而上） */
  background:
    linear-gradient(90deg, rgba(168,85,247,0.9), rgba(168,85,247,0.9)) bottom right / 0% 2px no-repeat,
    linear-gradient(180deg, rgba(124,58,237,0.9), rgba(124,58,237,0.9)) bottom left / 2px 0% no-repeat;
}

.media-card:hover::before,
.media-card:hover::after { opacity: 1; }
.media-card:hover::before { animation: edgesTopRight 0.6s linear forwards; }
.media-card:hover::after { animation: edgesBottomLeft 0.6s linear forwards; }

@keyframes edgesTopRight {
  0% { background-size: 0% 2px, 2px 0%; }
  100% { background-size: 100% 2px, 2px 100%; }
}

@keyframes edgesBottomLeft {
  0% { background-size: 0% 2px, 2px 0%; }
  100% { background-size: 100% 2px, 2px 100%; }
}

.poster {
  width: 100%;
  height: 70%;
  background-size: cover;
  background-position: center;
  position: relative;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.media-card:hover .poster {
  transform: scale(1.05);
}

/* 遮罩层 */
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.95) 100%);
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  opacity: 0;
  transition: opacity 0.4s ease;
}

.media-card:hover .overlay {
  opacity: 1;
}

.overlay-content {
  width: 90%;
  height: 90%;
  overflow-y: auto;
  padding: 30px 20px;
  text-align: left;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.3) transparent;
}

.overlay-content::-webkit-scrollbar {
  width: 6px;
}

.overlay-content::-webkit-scrollbar-track {
  background: transparent;
}

.overlay-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.overlay-content h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  border-bottom: 2px solid rgba(255, 255, 255, 0.3);
  padding-bottom: 10px;
}

.rating-display {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.stars { font-size: 1.1rem; line-height: 1; display: inline-flex; align-items: center; }
.star-wrap { position: relative; display: inline-block; width: 1em; height: 1em; margin-right: 2px; }
.star-wrap svg { width: 100%; height: 100%; display: block; }
/* 满星：实心黄色 - Font Awesome 实心星星 */
.stars .full-star,
.stars .full-star svg,
.stars .full-star svg path {
  color: #ffd700 !important;
}
/* 空星：黄色边框（空心）- Font Awesome 空心星星（far star 自动是空心的） */
.stars .empty-star,
.stars .empty-star svg,
.stars .empty-star svg path {
  color: #ffd700 !important;
}
/* 弹框中的星星样式 */
.modal-stars .full-star,
.modal-stars .full-star svg,
.modal-stars .full-star svg path {
  color: #ffd700 !important;
}
.modal-stars .empty-star,
.modal-stars .empty-star svg,
.modal-stars .empty-star svg path {
  color: #ffd700 !important;
}
.modal-stars .star-half-overlay,
.modal-stars .star-half-overlay svg,
.modal-stars .star-half-overlay svg path {
  color: #ffd700 !important;
}
.star-half-wrap { position: relative; display: block; width: 100%; height: 100%; }
.star-half-wrap .empty-star, .star-half-wrap .full-star { position: absolute; left: 0; top: 0; width: 100%; height: 100%; }
.star-half-overlay { position: absolute; left: 0; top: 0; width: 100%; height: 100%; clip-path: inset(0 50% 0 0); }

.rating-text { font-size: 0.9rem; font-weight: 600; color: #666; }

.overlay-content p {
  line-height: 1.8;
  font-size: 1rem;
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.edit-btn,
.delete-btn {
  flex: 1;
  padding: 10px 20px;
  border: none;
  border-radius: 25px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.edit-btn {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(168, 85, 247, 0.3);
}

.edit-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 25px rgba(168, 85, 247, 0.5);
}

.delete-btn {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(245, 87, 108, 0.3);
}

.delete-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 25px rgba(245, 87, 108, 0.5);
}

/* Markdown 渲染内容样式 */
.review-content.markdown-body {
  background: transparent !important;
  color: white !important;
  padding: 0;
  margin-bottom: 20px;
  max-height: 300px;
  overflow-y: auto;
}

.review-content.markdown-body h1,
.review-content.markdown-body h2,
.review-content.markdown-body h3,
.review-content.markdown-body h4,
.review-content.markdown-body h5,
.review-content.markdown-body h6 {
  color: white !important;
  border-bottom-color: rgba(255, 255, 255, 0.3) !important;
}

.review-content.markdown-body p,
.review-content.markdown-body li {
  color: rgba(255, 255, 255, 0.95) !important;
}

.review-content.markdown-body a {
  color: #a8d5ff !important;
}

.review-content.markdown-body code {
  background: rgba(255, 255, 255, 0.2) !important;
  color: #ffd700 !important;
}

.review-content.markdown-body pre {
  background: rgba(0, 0, 0, 0.3) !important;
}

.review-content.markdown-body blockquote {
  border-left-color: rgba(255, 255, 255, 0.5) !important;
  color: rgba(255, 255, 255, 0.9) !important;
}

/* 媒体模态框（对齐搜索模态视觉） */
.media-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 10000;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.media-modal-container {
  width: 90%;
  max-width: 600px; /* 与搜索模态一致 */
  max-height: 85vh;
  height: auto;
  min-height: 200px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  padding: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  cursor: default;
}

.media-modal-header {
  position: relative;
  min-height: 28px;
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: stretch;
  gap: 10px;
  padding-bottom: 0;
  margin-bottom: 0;
}

.modal-thumb { width: auto; max-width: none; border-radius: 8px; overflow: hidden; border: 1px solid rgba(0,0,0,0.08); background: #fff; display: flex; align-items: center; justify-content: center; }
.modal-thumb img { width: auto; max-width: none; height: 60px; max-height: 60px; object-fit: contain; display: block; }

.media-modal-headinfo { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; }
.media-modal-title { margin: 0; font-size: 1.1rem; font-weight: 700; color: #333; text-align: center; }
.media-modal-edit { justify-self: end; border: none; background: transparent; color: #a855f7; font-size: 0.9rem; cursor: pointer; padding: 2px 6px; }
.media-modal-edit:hover { color: #7c3aed; }

.media-modal-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  flex: 1;
  min-height: 100px;
  overflow: hidden;
  margin-top: 2.5rem;
}

.media-modal-content {
  flex: 1;
  min-width: 0;
  min-height: 100px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.media-modal-rating { display: flex; align-items: center; justify-content: center; gap: 8px; margin-bottom: 8px; }

.media-modal-review {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  min-height: 50px;
  text-align: left;
  padding-right: 8px;
}
.media-modal-review.markdown-body { background: transparent !important; padding: 0; }
.media-modal-review.markdown-body h1,
.media-modal-review.markdown-body h2,
.media-modal-review.markdown-body h3,
.media-modal-review.markdown-body h4,
.media-modal-review.markdown-body h5,
.media-modal-review.markdown-body h6 { background: transparent !important; }
/* 媒体弹框中正文图片样式 - 宽度不超过容器80%，高度不超过400px，居中显示 */
.media-modal-container .media-modal-review.markdown-body img,
.media-modal-content .media-modal-review.markdown-body img,
.media-modal-review.markdown-body img,
.media-modal-review img {
  max-width: 80% !important;
  max-height: 400px !important;
  width: auto !important;
  height: auto !important;
  display: block !important;
  margin: 20px auto !important;
  border-radius: 8px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
  object-fit: contain !important;
  box-sizing: border-box !important;
}

@media (max-width: 768px) {
  .media-modal-container { width: 90%; max-width: 90%; }
}

/* 美化滚动条样式 */
.media-modal-review {
  scrollbar-width: thin;
  scrollbar-color: rgba(168, 85, 247, 0.3) transparent;
  -ms-overflow-style: -ms-autohiding-scrollbar;
}

.media-modal-review::-webkit-scrollbar {
  width: 6px;
}

.media-modal-review::-webkit-scrollbar-track {
  background: transparent;
  border-radius: 3px;
}

.media-modal-review::-webkit-scrollbar-thumb {
  background: rgba(168, 85, 247, 0.3);
  border-radius: 3px;
  transition: background 0.2s ease;
}

.media-modal-review::-webkit-scrollbar-thumb:hover {
  background: rgba(168, 85, 247, 0.5);
}

.media-modal-actions { display: flex; gap: 10px; margin-top: 12px; }

.media-modal-fade-enter-active,
.media-modal-fade-leave-active { transition: opacity 0.3s ease; }
.media-modal-fade-enter-from,
.media-modal-fade-leave-to { opacity: 0; }

/* 卡片内容区 */
.content {
  display: flex;
  flex-direction: column;
  padding: 20px;
  height: 30%;
  justify-content: center;
  align-items: center;
  background: linear-gradient(to bottom, transparent, rgba(255, 255, 255, 0.95));
}

.media-card.media-row .content {
  height: auto;
  padding: 0;
  align-items: flex-start;
  background: transparent;
  width: 100%;
}
.row-header { display: flex; align-items: center; justify-content: space-between; width: 100%; gap: 12px; }
.row-title { margin: 0; font-size: 1rem; line-height: 1.3; font-weight: 600; color: #333; display: flex; align-items: center; gap: 6px; }
.row-rating { display: flex; align-items: center; gap: 6px; font-size: 0.9rem; color: inherit; white-space: nowrap; }
.row-meta { margin-top: 6px; font-size: 0.85rem; color: #999; }

h2 {
  margin: 0 0 8px 0;
  font-size: 1.3em;
  font-weight: 600;
  color: #333;
  text-align: center;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.type-fa-icon { color: #a855f7; margin-right: 6px; }

.media-card.media-row h2 {
  font-size: 1rem;
  line-height: 1.3;
  text-align: left;
  margin-bottom: 4px;
}

.date {
  margin-top: 5px;
  color: #999;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 5px;
}

.date::before { content: none; }

.rating-stars {
  margin-top: 8px;
  font-size: 1rem;
}

.media-card.media-row .rating-stars {
  margin-top: 2px;
  font-size: 0.9rem;
}

.rating-stars span {
  display: inline-block;
  letter-spacing: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .media-card {
    height: 400px;
  }

  .overlay-content {
    padding: 20px 15px;
  }
}
</style>
