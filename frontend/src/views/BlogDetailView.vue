<template>
  <div>
    <NavBar />
    <DetailArticleView :type="articleType" :article-id="$route.params.id" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '@/components/NavBar'
import DetailArticleView from '@/components/ArticleDetail.vue'

const route = useRoute()
// 从 query 参数中读取类型，如果没有则默认为 'blog'
const articleType = computed(() => route.query.type || 'blog')
</script>
