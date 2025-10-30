<template>
  <ArticleList
    type="all" class="basic-show-view"
    :headerImage="blogBackground"
    typingText="欢迎来到博客页面"
  />
</template>

<script setup>
import ArticleList from '@/components/ArticleList.vue'
import booksBackgroundImg from '@/assets/books_background.jpg'

const blogBackground = booksBackgroundImg
</script>
