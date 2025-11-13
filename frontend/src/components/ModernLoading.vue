<template>
  <div class="modern-loading-container">
    <!-- 背景动画已移除，使用主页背景 -->

    <!-- 主要内容 -->
    <div class="loading-content">
      <!-- Logo/头像区域 -->
      <div class="loading-avatar">
        <div class="avatar-container">
          <img src="@/assets/my_headportrait.jpg" alt="Avatar" class="avatar-image" />
          <div class="avatar-ring"></div>
          <div class="avatar-pulse"></div>
        </div>
      </div>

      <!-- 加载动画 -->
      <div class="loading-animation">
        <div class="spinner-container">
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
          <div class="spinner-ring"></div>
        </div>
      </div>

      <!-- 加载文本 -->
      <div class="loading-text">
        <h2 class="loading-title">{{ title }}</h2>
        <p class="loading-subtitle">{{ subtitle }}</p>
        <div class="loading-dots">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>

      <!-- 进度条 -->
      <div v-if="showProgress" class="loading-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: Math.min(100, Math.max(0, progress)) + '%' }"></div>
        </div>
        <div class="progress-text">{{ Math.min(100, Math.max(0, Math.round(progress))) }}%</div>
      </div>
    </div>

    <!-- 底部装饰 -->
    <div class="loading-footer">
      <div class="loading-tips">
        <div v-for="(tip, index) in tips" :key="index" class="tip-item" :class="{ active: currentTipIndex === index }">
          {{ tip }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

defineProps({
  title: {
    type: String,
    default: '山角函兽的小窝'
  },
  subtitle: {
    type: String,
    default: 'Loading……'
  },
  showProgress: {
    type: Boolean,
    default: true
  },
  progress: {
    type: Number,
    default: 0
  }
})

const currentTipIndex = ref(0)
const tips = ref([
  'Preparing content...',
  'Optimizing experience...',
  'Loading articles...',
  'Preparing recommendations...',
  'Almost ready...'
])

let tipInterval = null

onMounted(() => {
  // 轮播提示文本
  tipInterval = setInterval(() => {
    currentTipIndex.value = (currentTipIndex.value + 1) % tips.value.length
  }, 2000)
})

onUnmounted(() => {
  if (tipInterval) {
    clearInterval(tipInterval)
  }
})
</script>

<style scoped>
.modern-loading-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  overflow: hidden;
}

/* 使用现有的背景图片并添加模糊效果 */
.modern-loading-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('@/assets/background-image.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  filter: blur(6px);
  z-index: -2;
  height: 100%;
  width: 100vw;
}

/* 添加白蒙蒙的叠加层效果 */
.modern-loading-container::after {
  content: '';
  position: fixed;
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
  z-index: -1;
  height: 100vh;
  width: 100vw;
}

/* 背景动画已移除，使用主页背景 */

/* 主要内容 */
.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  z-index: 2;
  position: relative;
}

/* 头像区域 */
.loading-avatar {
  margin-bottom: 40px;
}

.avatar-container {
  position: relative;
  width: 120px;
  height: 120px;
}

.avatar-image {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  animation: avatarPulse 2s ease-in-out infinite;
}

.avatar-ring {
  position: absolute;
  top: -8px;
  left: -8px;
  right: -8px;
  bottom: -8px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
  animation: ringRotate 3s linear infinite;
}

.avatar-pulse {
  position: absolute;
  top: -12px;
  left: -12px;
  right: -12px;
  bottom: -12px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes avatarPulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes ringRotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.8;
  }
}

/* 加载动画 */
.loading-animation {
  margin-bottom: 30px;
}

.spinner-container {
  position: relative;
  width: 60px;
  height: 60px;
}

.spinner-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 3px solid transparent;
  border-top: 3px solid rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: spin 1.5s linear infinite;
}

.spinner-ring:nth-child(2) {
  width: 80%;
  height: 80%;
  top: 10%;
  left: 10%;
  border-top-color: rgba(255, 255, 255, 0.6);
  animation-duration: 1.2s;
  animation-direction: reverse;
}

.spinner-ring:nth-child(3) {
  width: 60%;
  height: 60%;
  top: 20%;
  left: 20%;
  border-top-color: rgba(255, 255, 255, 0.4);
  animation-duration: 0.8s;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* 文本区域 */
.loading-text {
  margin-bottom: 30px;
}

.loading-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 10px;
  text-shadow: 0 2px 10px rgba(255, 255, 255, 0.8);
  letter-spacing: 2px;
  animation: titleGlow 2s ease-in-out infinite alternate;
}

.loading-subtitle {
  font-size: 1.1rem;
  color: #666;
  margin-bottom: 20px;
  letter-spacing: 1px;
}

.loading-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.loading-dots span {
  width: 8px;
  height: 8px;
  background: #667eea;
  border-radius: 50%;
  animation: dotBounce 1.4s ease-in-out infinite both;
}

.loading-dots span:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-dots span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes titleGlow {
  0% {
    text-shadow: 0 2px 10px rgba(255, 255, 255, 0.8);
  }
  100% {
    text-shadow: 0 2px 10px rgba(255, 255, 255, 0.8), 0 0 20px rgba(102, 126, 234, 0.3);
  }
}

@keyframes dotBounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

/* 进度条 */
.loading-progress {
  width: 300px;
  margin-bottom: 30px;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 10px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.6));
  border-radius: 2px;
  transition: width 0.3s ease;
  position: relative;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  animation: progressShimmer 2s infinite;
}

@keyframes progressShimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.progress-text {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
  font-weight: 500;
}

/* 底部提示 */
.loading-footer {
  position: absolute;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2;
}

.loading-tips {
  text-align: center;
}

.tip-item {
  color: #888;
  font-size: 0.9rem;
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.5s ease;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}

.tip-item.active {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .loading-title {
    font-size: 2rem;
  }

  .loading-subtitle {
    font-size: 1rem;
  }

  .avatar-container {
    width: 100px;
    height: 100px;
  }

  .loading-progress {
    width: 250px;
  }
}

@media (max-width: 480px) {
  .loading-title {
    font-size: 1.5rem;
  }

  .loading-subtitle {
    font-size: 0.9rem;
  }

  .avatar-container {
    width: 80px;
    height: 80px;
  }

  .loading-progress {
    width: 200px;
  }
}
</style>
