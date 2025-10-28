<template>
  <ArticleList
    type="moment"
    class="basic-show-view"
    :headerImage="momentsBackground"
    typingText="记录生活的点点滴滴"
  />
</template>

<script setup>
import ArticleList from '@/components/ArticleList.vue'
import momentsBackgroundImg from '@/assets/moments_background.webp'

const momentsBackground = momentsBackgroundImg
</script>
