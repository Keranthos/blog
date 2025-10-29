<template>
  <div class="not-found-view">
    <NavBar />
    <div class="error-container">
      <div class="error-content">
        <!-- 现代化的404数字设计 -->
        <div class="error-code-container">
          <div class="error-code">404</div>
          <div class="error-code-shadow">404</div>
        </div>

        <!-- 现代化的标题和描述 -->
        <div class="error-title-section">
          <h1 class="error-title">页面走丢了...</h1>
          <div class="error-subtitle">Oops! Page Not Found</div>
        </div>

        <div class="error-message">
          <p>哎呀！你要找的页面好像不存在呢～</p>
          <p>可能是链接失效了，或者页面被移动到其他地方了。</p>
          <p>不过别担心，我们可以帮你找到回家的路！</p>
        </div>

        <!-- 现代化的操作按钮 -->
        <div class="error-actions">
          <button class="btn-primary" @click="goHome">
            <font-awesome-icon icon="home" />
            <span>回到首页</span>
            <div class="btn-shine"></div>
          </button>
          <button class="btn-secondary" @click="goBack">
            <font-awesome-icon icon="arrow-left" />
            <span>返回上一页</span>
            <div class="btn-shine"></div>
          </button>
        </div>

        <!-- 现代化的插图区域 -->
        <div class="error-illustration">
          <div class="illustration-container">
            <div class="floating-elements">
              <div class="floating-element element-1">🔍</div>
              <div class="floating-element element-2">📄</div>
              <div class="floating-element element-3">🌐</div>
              <div class="floating-element element-4">💫</div>
            </div>
            <div class="main-illustration">
              <div class="illustration-icon">🚀</div>
            </div>
          </div>
        </div>

        <!-- 现代化的装饰元素 -->
        <div class="decorative-elements">
          <div class="decorative-circle circle-1"></div>
          <div class="decorative-circle circle-2"></div>
          <div class="decorative-circle circle-3"></div>
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
  position: relative;
  background: transparent;
}

.not-found-view::before {
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

.not-found-view::after {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.8) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.4) 60%,
    rgba(255, 255, 255, 0.2) 100%
  );
  z-index: -1;
  height: 100vh;
  width: 100vw;
}

.error-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 80px);
  padding: 2rem;
  padding-top: 100px;
  position: relative;
  z-index: 2;
}

.error-content {
  text-align: center;
  color: #333;
  max-width: 800px;
  position: relative;
  z-index: 3;
}

/* 现代化的404数字设计 */
.error-code-container {
  position: relative;
  margin-bottom: 2rem;
}

.error-code {
  font-size: 12rem;
  font-weight: 900;
  line-height: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #5a67d8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 0 30px rgba(102, 126, 234, 0.3);
  animation: float 4s ease-in-out infinite;
  position: relative;
  z-index: 2;
}

.error-code-shadow {
  position: absolute;
  top: 0;
  left: 0;
  font-size: 12rem;
  font-weight: 900;
  line-height: 1;
  color: rgba(0, 0, 0, 0.1);
  z-index: 1;
  transform: translate(8px, 8px);
  animation: floatShadow 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-20px) scale(1.05);
  }
}

@keyframes floatShadow {
  0%, 100% {
    transform: translate(8px, 8px) scale(1);
  }
  50% {
    transform: translate(12px, 12px) scale(1.05);
  }
}

/* 现代化的标题区域 */
.error-title-section {
  margin-bottom: 2rem;
}

.error-title {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 4px 20px rgba(102, 126, 234, 0.3);
  letter-spacing: 2px;
  animation: titleGlow 3s ease-in-out infinite;
}

.error-subtitle {
  font-size: 1.2rem;
  font-weight: 400;
  opacity: 0.7;
  letter-spacing: 3px;
  text-transform: uppercase;
  margin-bottom: 1rem;
  color: #666;
}

@keyframes titleGlow {
  0%, 100% {
    text-shadow: 0 4px 20px rgba(255, 255, 255, 0.3);
  }
  50% {
    text-shadow: 0 4px 30px rgba(255, 255, 255, 0.6), 0 0 40px rgba(255, 255, 255, 0.3);
  }
}

.error-message {
  font-size: 1.1rem;
  margin-bottom: 3rem;
  line-height: 1.8;
  opacity: 0.9;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
  color: #555;
}

.error-message p {
  margin: 0.8rem 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

/* 现代化的操作按钮 */
.error-actions {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  margin-bottom: 3rem;
  flex-wrap: wrap;
}

.error-actions button {
  padding: 1rem 2.5rem;
  font-size: 1rem;
  font-weight: 600;
  border: none;
  border-radius: 50px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 0.8rem;
  position: relative;
  overflow: hidden;
  min-width: 180px;
  justify-content: center;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
  border: 2px solid rgba(102, 126, 234, 0.8);
}

.btn-primary:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 15px 35px rgba(102, 126, 234, 0.4);
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
}

.btn-secondary {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  backdrop-filter: blur(15px);
  border: 2px solid rgba(102, 126, 234, 0.3);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.1);
}

.btn-secondary:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 15px 35px rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.5);
}

.btn-shine {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.6s ease;
}

.error-actions button:hover .btn-shine {
  left: 100%;
}

/* 现代化的插图区域 */
.error-illustration {
  margin-top: 3rem;
  position: relative;
}

.illustration-container {
  position: relative;
  width: 300px;
  height: 300px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-illustration {
  width: 150px;
  height: 150px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  border: 3px solid rgba(102, 126, 234, 0.3);
  box-shadow: 0 20px 40px rgba(102, 126, 234, 0.2);
  animation: illustrationPulse 3s ease-in-out infinite;
}

.illustration-icon {
  font-size: 4rem;
  animation: iconBounce 2s ease-in-out infinite;
}

@keyframes illustrationPulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 20px 40px rgba(102, 126, 234, 0.2);
  }
  50% {
    transform: scale(1.1);
    box-shadow: 0 25px 50px rgba(102, 126, 234, 0.3);
  }
}

@keyframes iconBounce {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-10px) rotate(5deg);
  }
}

/* 浮动装饰元素 */
.floating-elements {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
}

.floating-element {
  position: absolute;
  font-size: 2rem;
  animation: floatAround 6s ease-in-out infinite;
  opacity: 0.7;
}

.element-1 {
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.element-2 {
  top: 30%;
  right: 15%;
  animation-delay: 1.5s;
}

.element-3 {
  bottom: 25%;
  left: 20%;
  animation-delay: 3s;
}

.element-4 {
  bottom: 20%;
  right: 10%;
  animation-delay: 4.5s;
}

@keyframes floatAround {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
    opacity: 0.7;
  }
  25% {
    transform: translateY(-20px) rotate(90deg);
    opacity: 1;
  }
  50% {
    transform: translateY(-10px) rotate(180deg);
    opacity: 0.8;
  }
  75% {
    transform: translateY(-15px) rotate(270deg);
    opacity: 0.9;
  }
}

/* 装饰圆圈 */
.decorative-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.decorative-circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.05) 100%);
  animation: circleFloat 8s ease-in-out infinite;
}

.circle-1 {
  width: 100px;
  height: 100px;
  top: 10%;
  left: 5%;
  animation-delay: 0s;
}

.circle-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 5%;
  animation-delay: 2s;
}

.circle-3 {
  width: 80px;
  height: 80px;
  bottom: 10%;
  left: 15%;
  animation-delay: 4s;
}

@keyframes circleFloat {
  0%, 100% {
    transform: translateY(0) scale(1);
    opacity: 0.3;
  }
  50% {
    transform: translateY(-30px) scale(1.1);
    opacity: 0.6;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .error-container {
    padding: 1rem;
  }

  .error-content {
    max-width: 100%;
  }

  .error-code {
    font-size: 8rem;
  }

  .error-code-shadow {
    font-size: 8rem;
  }

  .error-title {
    font-size: 2.2rem;
  }

  .error-subtitle {
    font-size: 1rem;
    letter-spacing: 2px;
  }

  .error-message {
    font-size: 1rem;
    margin-bottom: 2rem;
  }

  .error-actions {
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .error-actions button {
    width: 100%;
    min-width: auto;
    padding: 0.8rem 2rem;
  }

  .illustration-container {
    width: 250px;
    height: 250px;
  }

  .main-illustration {
    width: 120px;
    height: 120px;
  }

  .illustration-icon {
    font-size: 3rem;
  }

  .floating-element {
    font-size: 1.5rem;
  }

  .decorative-circle {
    opacity: 0.2;
  }

  .circle-1 {
    width: 60px;
    height: 60px;
  }

  .circle-2 {
    width: 80px;
    height: 80px;
  }

  .circle-3 {
    width: 50px;
    height: 50px;
  }
}

@media (max-width: 480px) {
  .error-code {
    font-size: 6rem;
  }

  .error-code-shadow {
    font-size: 6rem;
  }

  .error-title {
    font-size: 1.8rem;
  }

  .error-subtitle {
    font-size: 0.9rem;
    letter-spacing: 1px;
  }

  .error-message {
    font-size: 0.95rem;
    padding: 0 1rem;
  }

  .error-actions button {
    padding: 0.7rem 1.5rem;
    font-size: 0.9rem;
  }

  .illustration-container {
    width: 200px;
    height: 200px;
  }

  .main-illustration {
    width: 100px;
    height: 100px;
  }

  .illustration-icon {
    font-size: 2.5rem;
  }

  .floating-element {
    font-size: 1.2rem;
  }
}
</style>
