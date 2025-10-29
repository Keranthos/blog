<template>
  <transition name="fade">
    <div v-if="show" class="back-to-top-container">
      <!-- 滚动进度条 -->
      <div class="scroll-progress">
        <div
          class="progress-bar"
          :style="{ height: scrollProgress + '%' }"
        ></div>
      </div>

      <!-- 返回顶部按钮 -->
      <button
        class="back-to-top"
        :title="'返回顶部'"
        @click="scrollToTop"
      >
        <font-awesome-icon icon="arrow-up" />
      </button>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const show = ref(false)
const scrollProgress = ref(0)

// 滚动监听
const handleScroll = () => {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight

  // 计算滚动进度百分比
  scrollProgress.value = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0

  // 当滚动超过300px时显示按钮
  show.value = scrollTop > 300
}

// 滚动到顶部
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.back-to-top-container {
  position: fixed;
  bottom: 30px;
  left: 30px;
  z-index: 1000;
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.scroll-progress {
  width: 4px;
  height: 60px;
  background: rgba(102, 126, 234, 0.2);
  border-radius: 2px;
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.progress-bar {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  transition: height 0.1s ease;
  min-height: 2px;
}

.back-to-top {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(10px);
}

.back-to-top:hover {
  transform: translateY(-3px) scale(1.1);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.6);
  background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
}

.back-to-top:active {
  transform: translateY(-1px) scale(1.05);
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.8);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .back-to-top-container {
    bottom: 20px;
    left: 20px;
    gap: 10px;
  }

  .scroll-progress {
    height: 50px;
  }

  .back-to-top {
    width: 45px;
    height: 45px;
    font-size: 1.1rem;
  }
}

@media (max-width: 480px) {
  .back-to-top-container {
    bottom: 15px;
    left: 15px;
    gap: 8px;
  }

  .scroll-progress {
    height: 40px;
    width: 3px;
  }

  .back-to-top {
    width: 40px;
    height: 40px;
    font-size: 1rem;
  }
}

/* 滚动进度条样式 */
.back-to-top::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border-radius: 50%;
  background: conic-gradient(
    from 0deg,
    #667eea 0deg,
    #764ba2 360deg
  );
  z-index: -1;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.back-to-top:hover::before {
  opacity: 0.3;
  animation: rotate 2s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
