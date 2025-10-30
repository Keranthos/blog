<template>
  <div class="media-card" @mouseover="showOverlay = true" @mouseleave="showOverlay = false">
    <div class="poster" :style="{ backgroundImage: `url(${props.media.Poster})` }"></div>
    <div class="content">
      <h2>{{ props.media.Name }}</h2>
      <p class="date">创建日期: {{ formattedTime }}</p>
      <div v-if="props.media.Rating" class="rating-stars">
        <span>{{ '⭐'.repeat(Math.round(props.media.Rating / 2)) }}</span>
      </div>
    </div>
    <!-- 黑色背景遮罩层 -->
    <div v-if="showOverlay" class="overlay">
      <div class="overlay-content">
        <h3>{{ props.media.Name }}</h3>
        <div v-if="props.media.Rating" class="rating-display">
          <span class="stars">{{ '⭐'.repeat(Math.round(props.media.Rating / 2)) }}</span>
          <span class="rating-text">{{ props.media.Rating }}/10</span>
        </div>
        <!-- 使用 Markdown 渲染 -->
        <div class="review-content markdown-body" v-html="renderedReview"></div>
        <div v-if="userLevel >= 3" class="action-buttons">
          <button class="edit-btn" @click="goToEdit">
            <font-awesome-icon icon="edit" /> 编辑
          </button>
          <button class="delete-btn" @click="confirmDelete">
            <font-awesome-icon icon="trash" /> 删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineEmits, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { deleteMedia } from '@/api/media/edit'

const props = defineProps({
  media: Object,
  type: String
})

const emit = defineEmits(['mediaDeleted'])

const store = useStore()
const router = useRouter()
const user = store.state.user
const userLevel = user.level

const showOverlay = ref(false)

// 格式化时间
const formattedTime = computed(() => {
  const date = new Date(props.media.CreatedAt)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
})

// 渲染 Markdown
const renderedReview = computed(() => {
  if (!props.media.Review) return ''
  const html = marked(props.media.Review)
  return DOMPurify.sanitize(html)
})

// 跳转到编辑页面
const goToEdit = () => {
  router.push({
    path: `/edit/${props.media.ID}`,
    query: {
      contentType: 'media',
      mediaType: props.type
    }
  })
}

const confirmDelete = async () => {
  // 第一次确认
  const firstConfirm = confirm(`确定要删除《${props.media.Name}》吗？\n\n⚠️ 此操作不可撤销！`)
  if (!firstConfirm) return

  // 第二次确认（防误触）
  const secondConfirm = confirm(`再次确认：真的要删除《${props.media.Name}》吗？\n\n删除后将无法恢复！`)
  if (!secondConfirm) return

  try {
    await deleteMedia(props.type, props.media.ID)
    emit('mediaDeleted', props.media.ID)

    // 显示看板娘消息（如果可用）
    if (window.showMessage) {
      window.showMessage('删除成功～', 3000, 9)
    }
  } catch (error) {
    console.error('Failed to delete media:', error)

    // 显示看板娘错误消息（如果可用）
    if (window.showMessage) {
      window.showMessage('(｡•́︿•̀｡)<br>删除失败了…请重试吧～', 5000, 10)
    }
  }
}
</script>

<style>
/* 导入 Markdown 样式（全局） */
@import '@/assets/styles/github-markdown.css';
</style>

<style scoped>
.media-card {
  position: relative;
  cursor: pointer;
  border-radius: 16px;
  overflow: hidden;
  background: white;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  height: 450px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.media-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
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

.stars {
  font-size: 1.2rem;
}

.rating-text {
  font-size: 1.1rem;
  font-weight: 600;
  color: #ffd700;
}

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

.date {
  margin-top: 5px;
  color: #999;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 5px;
}

.date::before {
  content: '📅';
}

.rating-stars {
  margin-top: 8px;
  font-size: 1rem;
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
