<template>
  <div class="loading-test-view">
    <NavBar />

    <div class="test-container">
      <h1>加载界面测试页面</h1>
      <p>点击下面的按钮来测试不同的加载效果</p>

      <div class="test-buttons">
        <button class="test-btn" @click="testHomeLoading">测试首页加载</button>
        <button class="test-btn" @click="testArticleLoading">测试文章列表加载</button>
        <button class="test-btn" @click="testCustomLoading">测试自定义加载</button>
      </div>

      <div class="test-info">
        <h3>测试说明：</h3>
        <ul>
          <li>首页加载：模拟首页初始化过程</li>
          <li>文章列表加载：模拟文章列表加载过程</li>
          <li>自定义加载：自定义标题和进度的加载</li>
        </ul>
      </div>
    </div>

    <!-- 加载界面 -->
    <ModernLoading
      v-if="showLoading"
      :progress="loadingProgress"
      :title="loadingTitle"
      :subtitle="loadingSubtitle"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import NavBar from '@/components/NavBar.vue'
import ModernLoading from '@/components/ModernLoading.vue'

const showLoading = ref(false)
const loadingProgress = ref(0)
const loadingTitle = ref('')
const loadingSubtitle = ref('')

const testHomeLoading = () => {
  showLoading.value = true
  loadingProgress.value = 0
  loadingTitle.value = '山角函兽的小窝'
  loadingSubtitle.value = 'Loading……'

  simulateLoading()
}

const testArticleLoading = () => {
  showLoading.value = true
  loadingProgress.value = 0
  loadingTitle.value = '我的博客'
  loadingSubtitle.value = 'Loading……'

  simulateLoading()
}

const testCustomLoading = () => {
  showLoading.value = true
  loadingProgress.value = 0
  loadingTitle.value = '自定义加载测试'
  loadingSubtitle.value = 'Loading……'

  simulateLoading()
}

const simulateLoading = () => {
  const progressInterval = setInterval(() => {
    if (loadingProgress.value < 90) {
      loadingProgress.value += Math.random() * 15
    }
  }, 200)

  // 模拟加载过程
  setTimeout(() => {
    loadingProgress.value = 100
    clearInterval(progressInterval)

    setTimeout(() => {
      showLoading.value = false
      loadingProgress.value = 0
    }, 500)
  }, 3000)
}
</script>

<style scoped>
.loading-test-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.test-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  text-align: center;
}

.test-container h1 {
  font-size: 2.5rem;
  margin-bottom: 20px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.test-container p {
  font-size: 1.2rem;
  margin-bottom: 40px;
  opacity: 0.9;
}

.test-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.test-btn {
  padding: 15px 30px;
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 25px;
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.test-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
}

.test-info {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  padding: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  text-align: left;
}

.test-info h3 {
  margin-bottom: 15px;
  color: #fff;
}

.test-info ul {
  list-style: none;
  padding: 0;
}

.test-info li {
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.test-info li:last-child {
  border-bottom: none;
}

.test-info li::before {
  content: '✨';
  margin-right: 10px;
}

@media (max-width: 768px) {
  .test-container h1 {
    font-size: 2rem;
  }

  .test-buttons {
    flex-direction: column;
    align-items: center;
  }

  .test-btn {
    width: 200px;
  }
}
</style>
