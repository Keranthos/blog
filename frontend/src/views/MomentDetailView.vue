<template>
  <div>
    <NavBar />
    <ArticleDetail :articleId="articleId" :type="'moment'" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import ArticleDetail from '@/components/ArticleDetail.vue'

const route = useRoute()
const articleId = computed(() => route.params.id)
</script>
