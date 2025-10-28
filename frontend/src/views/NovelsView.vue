<template>
  <div class="media-view">
    <NavBar />
    <div class="header">
      <img :src="novelsBackground" alt="Background Image" class="header-image" loading="lazy" />
    </div>
    <div class="message1">
      <p>这是我的小说</p>
    </div>
    <div class="media-list" @scroll="handleScroll">
      <MediaCard v-for="media in mediaList" :key="media.ID" :media="media" type="novels" @media-deleted="removeMedia" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, defineAsyncComponent } from 'vue'
import { getMediaList } from '@/api/media/browse'
import MediaCard from '@/components/MediaCard.vue'
import novelsBackgroundImg from '@/assets/novels_background.jpg'

const NavBar = defineAsyncComponent(() => import('@/components/NavBar.vue'))

const novelsBackground = novelsBackgroundImg
const mediaList = ref([])
const currentPage = ref(1)
const isLoading = ref(false)
const mediaCache = ref({})

const loadMedia = async () => {
  if (isLoading.value) return
  isLoading.value = true
  try {
    if (mediaCache.value[currentPage.value]) {
      mediaList.value.push(...mediaCache.value[currentPage.value])
    } else {
      const response = await getMediaList('novels', currentPage.value)
      mediaList.value.push(...response.data.medias)
      mediaCache.value[currentPage.value] = response.data.medias
    }
    currentPage.value++
  } catch (error) {
    console.error('Failed to load media:', error)
  } finally {
    isLoading.value = false
  }
}

const handleScroll = (event) => {
  const { scrollTop, scrollHeight, clientHeight } = event.target
  if (scrollTop + clientHeight >= scrollHeight - 10) {
    loadMedia()
  }
}

const removeMedia = (mediaId) => {
  mediaList.value = mediaList.value.filter(media => media.ID !== mediaId)
  for (const page in mediaCache.value) {
    mediaCache.value[page] = mediaCache.value[page].filter(media => media.ID !== mediaId)
  }
}

onMounted(() => {
  loadMedia()
})
</script>

<style scoped>
.media-view {
  min-height: 100vh;
  padding-top: 40px;
  padding-bottom: 80px;
}

.header {
  position: relative;
  width: 100%;
  padding: 40px 350px;
  margin-bottom: 40px;
  overflow: hidden;
}

.header-image {
  width: 100%;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  display: block;
  margin: 0 auto;
}

.message1 {
  text-align: center;
  margin: 40px 0;
  color: white;
  font-size: 2.5rem;
  font-weight: 700;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  letter-spacing: 2px;
}

.media-list {
  width: 90%;
  max-width: 1400px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 30px;
  padding-bottom: 60px;
}

.create-media {
  width: 90%;
  max-width: 800px;
  margin: 40px auto 0;
  padding: 30px;
  background: transparent;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
}

.create-media h3 {
  color: #333;
  margin-bottom: 20px;
  font-size: 1.5rem;
}

.create-media input,
.create-media textarea {
  width: 100%;
  padding: 15px;
  margin-bottom: 15px;
  border: 2px solid #e8eaf0;
  border-radius: 12px;
  font-size: 1rem;
  font-family: inherit;
  transition: all 0.3s ease;
}

.create-media input:focus,
.create-media textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.create-media textarea {
  min-height: 120px;
  resize: vertical;
}

.create-media button {
  padding: 12px 30px;
  border: none;
  border-radius: 25px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.create-media button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
}

@media (max-width: 1400px) {
  .media-list {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 1024px) {
  .media-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .media-view {
    padding-top: 80px;
  }

  .header-image {
    height: 250px;
  }

  .message1 {
    font-size: 1.8rem;
  }

  .media-list {
    width: 95%;
    grid-template-columns: 1fr;
  }

  .create-media {
    width: 95%;
    padding: 20px;
  }
}
</style>
