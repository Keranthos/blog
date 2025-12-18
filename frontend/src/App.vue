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
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { resetSEO } from '@/utils/seo'
import BackToTop from '@/components/BackToTop.vue'

const route = useRoute()

// 初始化
onMounted(() => {
  // 重置SEO信息
  resetSEO()
})

// 监听路由变化，确保标题及时更新
watch(() => route.path, (newPath, oldPath) => {
  // 路由变化时，如果不是文章详情页，重置SEO
  // 文章详情页会在自己的组件中更新SEO
  // 使用更精确的匹配，避免误判
  const isArticleDetailPage = /^\/blog\/\d+/.test(newPath) || /^\/moments\/\d+/.test(newPath)

  if (!isArticleDetailPage && oldPath) {
    // 只有在真正切换路由时才重置（oldPath存在说明不是首次加载）
    // 延迟一点确保路由完全切换，但不要太长，避免与组件内的SEO更新冲突
    setTimeout(() => {
      resetSEO()
    }, 50)
  }
}, { immediate: false })
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

/* 确保 Emoji 正确显示 */
.comment-content,
.comment-text,
.comment-bubble,
.comment-input,
textarea {
  font-family: 'Inter', 'Noto Sans SC', 'Apple Color Emoji', 'Segoe UI Emoji', 'Noto Color Emoji', 'Segoe UI Symbol', 'Android Emoji', 'EmojiSymbols', 'EmojiOne Mozilla', 'Twemoji Mozilla', Avenir, Helvetica, Arial, sans-serif;
}

/* 确保页面内容从顶部开始 */
body, html {
  margin: 0;
  padding: 0;
  min-height: 100vh;
}
</style>
