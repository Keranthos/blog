<template>
  <Teleport to="body">
    <Transition name="share-card-fade">
      <div v-if="visible" class="share-card-overlay" @click="close">
        <div class="share-card-container" @click.stop>
          <!-- 顶部品牌区 -->
          <div class="share-card-header">
            <div class="brand-icon">
              <img :src="require('@/assets/my_headportrait.jpg')" alt="Avatar" />
            </div>
            <div class="brand-info">
              <div class="brand-name">山角函兽</div>
            </div>
          </div>
          <div class="header-divider"></div>

          <!-- 主要内容区 -->
          <div class="share-card-content">
            <div class="selected-text">
              {{ selectedText }}
            </div>
          </div>
          <div class="content-divider"></div>

          <!-- 底部信息区 -->
          <div class="share-card-footer">
            <div class="footer-content">
              <div class="footer-left">
                <div class="article-title">{{ articleTitle }}</div>
                <div class="copyright">Powered By 山角函兽 | 于 {{ currentDate }}</div>
              </div>
              <!-- 二维码（底部右侧） -->
              <div class="qr-code-wrapper">
                <canvas ref="qrCanvas" class="qr-code"></canvas>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch, nextTick, onMounted } from 'vue'
import QRCode from 'qrcode'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  selectedText: {
    type: String,
    default: ''
  },
  articleTitle: {
    type: String,
    default: ''
  },
  articleSubtitle: {
    type: String,
    default: ''
  },
  articleUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close'])

const qrCanvas = ref(null)
const currentDate = ref('')

// 生成当前日期
const generateDate = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  currentDate.value = `${year}/${month}/${day}`
}

// 生成二维码
const generateQRCode = async () => {
  if (!qrCanvas.value || !props.articleUrl) return

  try {
    await QRCode.toCanvas(qrCanvas.value, props.articleUrl, {
      width: 80,
      margin: 1,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    })
  } catch (error) {
    console.error('生成二维码失败:', error)
  }
}

// 监听可见性变化，生成二维码
watch(() => props.visible, async (newVal) => {
  if (newVal) {
    generateDate()
    await nextTick()
    generateQRCode()
  }
})

// 监听 URL 变化，重新生成二维码
watch(() => props.articleUrl, () => {
  if (props.visible) {
    generateQRCode()
  }
})

const close = () => {
  emit('close')
}

onMounted(() => {
  generateDate()
})
</script>

<style scoped>
.share-card-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
  backdrop-filter: blur(4px);
}

.share-card-container {
  position: relative;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  background-image: url('@/assets/background-image.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(20px);
  /* 添加白色半透明叠加层，类似主界面效果 */
  position: relative;
}

.share-card-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.7) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0.4) 100%
  );
  z-index: 0;
  pointer-events: none;
}

.share-card-container > * {
  position: relative;
  z-index: 1;
}

/* 顶部品牌区 */
.share-card-header {
  display: flex;
  align-items: center;
  padding: 16px 20px 12px 20px;
  background: transparent;
}

.brand-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 12px;
  flex-shrink: 0;
  border: 2px solid rgba(168, 85, 247, 0.3);
  box-shadow: 0 2px 8px rgba(168, 85, 247, 0.2);
}

.brand-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.brand-info {
  flex: 1;
}

.brand-name {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1.4;
  letter-spacing: 0.3px;
}

.header-divider {
  height: 1px;
  background: rgba(168, 85, 247, 0.3);
  margin: 0 20px;
}

/* 主要内容区 */
.share-card-content {
  flex: 1;
  padding: 24px 20px;
  position: relative;
  min-height: 180px;
  overflow-y: auto;
  background: transparent;
}

.selected-text {
  font-size: 18px;
  line-height: 2;
  color: #1a1a1a;
  text-align: left;
  margin-bottom: 0;
  word-break: break-word;
}

/* 选中文本中的高亮词（可选，用于某些关键词高亮） */
.selected-text :deep(mark) {
  background: rgba(168, 85, 247, 0.4);
  color: inherit;
  padding: 2px 4px;
  border-radius: 3px;
}

.content-divider {
  height: 2px;
  background: rgba(168, 85, 247, 0.6);
  margin: 0;
}

.qr-code {
  width: 100%;
  height: 100%;
}

/* 底部信息区 */
.share-card-footer {
  padding: 16px 20px;
  background: rgba(124, 58, 237, 0.95);
  color: white;
}

.footer-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.footer-left {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  line-height: 1.5;
  word-break: break-word;
}

.copyright {
  font-size: 13px;
  opacity: 0.9;
  margin-top: 8px;
  font-weight: 400;
}

/* 二维码（底部右侧） */
.qr-code-wrapper {
  width: 67px;
  height: 67px;
  background: white;
  border-radius: 8px;
  padding: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* 过渡动画 */
.share-card-fade-enter-active,
.share-card-fade-leave-active {
  transition: opacity 0.3s ease;
}

.share-card-fade-enter-active .share-card-container,
.share-card-fade-leave-active .share-card-container {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.share-card-fade-enter-from {
  opacity: 0;
}

.share-card-fade-leave-to {
  opacity: 0;
}

.share-card-fade-enter-from .share-card-container {
  transform: scale(0.9) translateY(20px);
  opacity: 0;
}

.share-card-fade-leave-to .share-card-container {
  transform: scale(0.95);
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .share-card-container {
    width: 95%;
    max-width: 95%;
    max-height: 85vh;
  }

  .footer-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .qr-code-wrapper {
    width: 53px;
    height: 53px;
    align-self: flex-end;
  }

  .share-card-header {
    padding: 16px 20px;
  }

  .share-card-content {
    padding: 24px 20px;
  }

  .share-card-footer {
    padding: 16px 20px;
  }

  .brand-name {
    font-size: 16px;
  }

  .selected-text {
    font-size: 16px;
  }

  .article-title {
    font-size: 16px;
  }

  .copyright {
    font-size: 12px;
  }
}
</style>
