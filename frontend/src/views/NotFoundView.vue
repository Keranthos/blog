<template>
  <div class="not-found-view">
    <NavBar />
    <div class="error-container">
      <div class="error-content">
        <div class="error-code">404</div>
        <div class="error-title">页面走丢了...</div>
        <div class="error-message">
          <p>哎呀！你要找的页面好像不存在呢～</p>
          <p>可能是链接失效了，或者页面被移动到其他地方了。</p>
        </div>
        <div class="error-actions">
          <button class="btn-primary" @click="goHome">
            <font-awesome-icon icon="home" />
            回到首页
          </button>
          <button class="btn-secondary" @click="goBack">
            <font-awesome-icon icon="arrow-left" />
            返回上一页
          </button>
        </div>
        <div class="error-illustration">
          <img src="https://picsum.photos/id/404/600/400" alt="404 Illustration" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped>
.not-found-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.error-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 60px);
  padding: 2rem;
}

.error-content {
  text-align: center;
  color: white;
  max-width: 600px;
}

.error-code {
  font-size: 10rem;
  font-weight: 900;
  line-height: 1;
  margin-bottom: 1rem;
  text-shadow: 4px 4px 8px rgba(0, 0, 0, 0.3);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}

.error-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
}

.error-message {
  font-size: 1.1rem;
  margin-bottom: 2.5rem;
  line-height: 1.6;
  opacity: 0.95;
}

.error-message p {
  margin: 0.5rem 0;
}

.error-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-bottom: 2rem;
}

.error-actions button {
  padding: 0.8rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  border: none;
  border-radius: 25px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-primary {
  background: white;
  color: #667eea;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  backdrop-filter: blur(10px);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.error-illustration {
  margin-top: 2rem;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.error-illustration img {
  width: 100%;
  height: auto;
  display: block;
}

@media (max-width: 768px) {
  .error-code {
    font-size: 6rem;
  }

  .error-title {
    font-size: 1.8rem;
  }

  .error-message {
    font-size: 1rem;
  }

  .error-actions {
    flex-direction: column;
  }

  .error-actions button {
    width: 100%;
  }
}
</style>
