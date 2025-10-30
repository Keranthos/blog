<template>
  <div class="media-card-wrapper">
  <div class="media-card" @click="openMediaModal">
    <div class="poster" :style="{ backgroundImage: `url(${props.media.Poster})` }"></div>
    <div class="content">
      <h2>{{ props.media.Name }}</h2>
      <p class="date">创建日期: {{ formattedTime }}</p>
      <div v-if="props.media.Rating" class="rating-stars">
        <span>{{ '⭐'.repeat(Math.round(props.media.Rating / 2)) }}</span>
      </div>
    </div>

  </div>

  <!-- 媒体详情模态框 - 与搜索模态风格一致 -->
  <Teleport to="body">
    <transition name="media-modal-fade">
      <div v-if="showMediaModal" class="media-modal-overlay" @click="closeMediaModal">
        <div class="media-modal-container" @click.stop>
          <div class="media-modal-header">
            <h3 class="media-modal-title">{{ props.media.Name }}</h3>
            <button v-if="userLevel >= 3" class="media-modal-edit" @click="goToEdit">编辑</button>
          </div>
          <div class="media-modal-body">
            <div class="media-modal-content">
              <div v-if="props.media.Rating" class="media-modal-rating">
                <span class="stars">{{ '⭐'.repeat(Math.round(props.media.Rating / 2)) }}</span>
                <span class="rating-text">{{ props.media.Rating }}/10</span>
              </div>
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
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps({
  media: Object,
  type: String
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

// 渲染 Markdown
const renderedReview = computed(() => {
  if (!props.media.Review) return ''
  const html = marked(props.media.Review)
  return DOMPurify.sanitize(html)
})

const openMediaModal = () => {
  showMediaModal.value = true
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
}

.media-modal-title { position: absolute; left: 50%; top: 0; transform: translateX(-50%); margin: 0; font-size: 1.1rem; font-weight: 700; color: #333; text-align: center; }
.media-modal-edit { position: absolute; right: 0; top: -2px; border: none; background: transparent; color: #a855f7; font-size: 0.9rem; cursor: pointer; padding: 2px 6px; }
.media-modal-edit:hover { color: #7c3aed; }

.media-modal-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
  min-height: 0;
}

.media-modal-content { flex: 1; min-width: 0; display: flex; flex-direction: column; }

.media-modal-rating { display: flex; align-items: center; justify-content: center; gap: 8px; margin-bottom: 8px; }

.media-modal-review { flex: 1; overflow-y: auto; min-height: 0; text-align: left; }
.media-modal-review.markdown-body { background: transparent !important; padding: 0; }
.media-modal-review.markdown-body h1,
.media-modal-review.markdown-body h2,
.media-modal-review.markdown-body h3,
.media-modal-review.markdown-body h4,
.media-modal-review.markdown-body h5,
.media-modal-review.markdown-body h6 { background: transparent !important; }

@media (max-width: 768px) {
  .media-modal-container { width: 90%; max-width: 90%; }
}
.media-modal-review { scrollbar-width: none; -ms-overflow-style: none; }
.media-modal-review::-webkit-scrollbar { display: none; }

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
