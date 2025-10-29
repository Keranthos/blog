<template>
  <div>
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component" />
      </keep-alive>
    </router-view>

    <!-- 全局返回顶部按钮 -->
    <BackToTop />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { resetSEO } from '@/utils/seo'
import BackToTop from '@/components/BackToTop.vue'

// 初始化
onMounted(() => {
  // 重置SEO信息
  resetSEO()
})
</script>

<style>
#app {
  font-family: 'Inter', 'Noto Sans SC', Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  min-height: 100vh;
}

/* 确保页面内容从顶部开始 */
body, html {
  margin: 0;
  padding: 0;
  min-height: 100vh;
}
</style>
