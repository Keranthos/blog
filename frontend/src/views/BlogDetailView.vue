<template>
  <div>
    <NavBar />
    <DetailArticleView type="blog" :article-id="$route.params.id" />
  </div>
</template>

<script setup>
import NavBar from '@/components/NavBar'
import DetailArticleView from '@/components/ArticleDetail.vue'
</script>
