<template>
  <ArticleList
    type="all" class="basic-show-view"
    :headerImage="blogBackground"
    typingText="欢迎来到博客页面"
  />
</template>

<script setup>
import ArticleList from '@/components/ArticleList.vue'
import blogBackgroundImg from '@/assets/blog_background.jpg'

const blogBackground = blogBackgroundImg
</script>
